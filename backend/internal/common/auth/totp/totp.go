package totp

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
)

type TOTPService struct {
	issuer string
}

type TOTPSetup struct {
	Secret      string   `json:"secret"`
	QRCode      string   `json:"qr_code"`
	URL         string   `json:"url"`
	BackupCodes []string `json:"backup_codes"`
}

type BackupCodePair struct {
	Plain  []string
	Hashed []string
}

func NewTOTPService(issuer string) *TOTPService {
	return &TOTPService{
		issuer: issuer,
	}
}

func (s *TOTPService) GenerateSecret(accountName string) (*TOTPSetup, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: accountName,
		Period:      30,
		SecretSize:  32,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	qrCode, err := s.generateQRCode(key.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	plainCodes, _, err := s.GenerateBackupCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}

	return &TOTPSetup{
		Secret:      key.Secret(),
		QRCode:      qrCode,
		URL:         key.URL(),
		BackupCodes: plainCodes,
	}, nil
}

func (s *TOTPService) generateQRCode(url string) (string, error) {
	qr, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(qr)
	return fmt.Sprintf("data:image/png;base64,%s", encoded), nil
}

func (s *TOTPService) ValidateToken(secret, token string) bool {
	valid, err := totp.ValidateCustom(
		token,
		secret,
		time.Now(),
		totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	if err != nil {
		return false
	}
	return valid
}

func (s *TOTPService) GenerateBackupCodes() ([]string, []string, error) {
	plainCodes := make([]string, 10)
	hashedCodes := make([]string, 10)

	for i := 0; i < 10; i++ {
		code, err := s.generateBackupCode()
		if err != nil {
			return nil, nil, err
		}
		plainCodes[i] = code

		hashed, err := bcrypt.GenerateFromPassword([]byte(code), 12)
		if err != nil {
			return nil, nil, err
		}
		hashedCodes[i] = string(hashed)
	}

	return plainCodes, hashedCodes, nil
}

func (s *TOTPService) generateBackupCode() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	code := fmt.Sprintf("%02d%02d-%02d%02d-%02d%02d",
		bytes[0]%100, bytes[1]%100,
		bytes[2]%100, bytes[3]%100,
		bytes[4]%100, bytes[5]%100)
	return code, nil
}

func (s *TOTPService) ValidateBackupCode(hashedCodesJSON string, inputCode string) (int, bool) {
	var hashedCodes []string
	if err := json.Unmarshal([]byte(hashedCodesJSON), &hashedCodes); err != nil {
		return -1, false
	}

	for i, hashedCode := range hashedCodes {
		if hashedCode == "" || hashedCode == "used" {
			continue
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(inputCode)); err == nil {
			return i, true
		}
	}

	return -1, false
}

func (s *TOTPService) MarkBackupCodeUsed(hashedCodesJSON string, index int) (string, error) {
	var hashedCodes []string
	if err := json.Unmarshal([]byte(hashedCodesJSON), &hashedCodes); err != nil {
		return "", err
	}

	if index < 0 || index >= len(hashedCodes) {
		return "", fmt.Errorf("invalid backup code index")
	}

	hashedCodes[index] = "used"

	updatedJSON, err := json.Marshal(hashedCodes)
	if err != nil {
		return "", err
	}

	return string(updatedJSON), nil
}

func (s *TOTPService) CountRemainingBackupCodes(hashedCodesJSON string) int {
	if hashedCodesJSON == "" {
		return 0
	}

	var hashedCodes []string
	if err := json.Unmarshal([]byte(hashedCodesJSON), &hashedCodes); err != nil {
		return 0
	}

	count := 0
	for _, code := range hashedCodes {
		if code != "" && code != "used" {
			count++
		}
	}

	return count
}
