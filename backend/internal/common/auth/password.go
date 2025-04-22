package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Cost for bcrypt hashing (12 is a good balance for security vs performance)
const BcryptCost = 12

// Common password validation errors
var (
	ErrPasswordTooShort = errors.New("password is too short, must be at least 8 characters")
	ErrPasswordTooWeak  = errors.New("password is too weak, must contain uppercase, lowercase, digit, and special character")
	ErrPasswordMismatch = errors.New("password does not match")
)

// HashPassword generates a bcrypt hash from a password
func HashPassword(password string) (string, error) {
	// Validate password strength before hashing
	if err := ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	// Generate bcrypt hash with appropriate cost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

// VerifyPassword checks if a password matches a bcrypt hash
func VerifyPassword(hashedPassword, password string) error {
	// Compare password with hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}

	return nil
}

// ValidatePasswordStrength checks if a password meets security requirements
func ValidatePasswordStrength(password string) error {
	// Check password length
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	// Check for at least one uppercase letter, one lowercase letter, one digit, and one special character
	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' ||
			char == '^' || char == '&' || char == '*' || char == '(' || char == ')' ||
			char == '-' || char == '_' || char == '+' || char == '=' || char == '{' ||
			char == '}' || char == '[' || char == ']' || char == '|' || char == '\\' ||
			char == ':' || char == ';' || char == '"' || char == '\'' || char == '<' ||
			char == '>' || char == ',' || char == '.' || char == '?' || char == '/':
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasDigit && hasSpecial) {
		return ErrPasswordTooWeak
	}

	return nil
}

// ComparePasswords checks if two passwords match (for password confirmation)
func ComparePasswords(password, confirmPassword string) error {
	if password != confirmPassword {
		return ErrPasswordMismatch
	}
	return nil
}
