package contact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type EmailConfig struct {
	ResendAPIKey string
	FromEmail    string
	ToEmail      string
}

type EmailMessage struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func (cfg *EmailConfig) IsConfigured() bool {
	return cfg.ResendAPIKey != "" && cfg.ToEmail != ""
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	ReplyTo string   `json:"reply_to,omitempty"`
}

func SendEmail(cfg *EmailConfig, msg *EmailMessage) error {
	if !cfg.IsConfigured() {
		return fmt.Errorf("email not configured")
	}

	subject := "Portfolio Contact: " + msg.Name
	if msg.Subject != "" {
		subject = "Portfolio Contact: " + msg.Name + " - " + msg.Subject
	}

	fromAddr := cfg.FromEmail
	if fromAddr == "" {
		fromAddr = "Contact Form <onboarding@resend.dev>"
	}

	body := buildEmailBody(msg)

	reqBody := resendRequest{
		From:    fromAddr,
		To:      []string{cfg.ToEmail},
		Subject: subject,
		Text:    body,
		ReplyTo: msg.Email,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ResendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func buildEmailBody(msg *EmailMessage) string {
	var b strings.Builder
	b.WriteString("New Contact Form Submission\n")
	b.WriteString("==========================\n\n")
	b.WriteString(fmt.Sprintf("Name:    %s\n", msg.Name))
	b.WriteString(fmt.Sprintf("Email:   %s\n", msg.Email))
	if msg.Subject != "" {
		b.WriteString(fmt.Sprintf("Subject: %s\n", msg.Subject))
	}
	b.WriteString(fmt.Sprintf("\nMessage:\n--------\n%s\n", msg.Message))
	return b.String()
}
