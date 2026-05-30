package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"echobackend/config"
)

const resendAPIURL = "https://api.resend.com/emails"

// Service sends application emails through Resend.
type Service struct {
	apiKey     string
	from       string
	httpClient *http.Client
}

// NewService creates a Resend-backed email service. An empty API key disables delivery.
func NewService(cfg config.EmailConfig) *Service {
	return &Service{
		apiKey: cfg.ResendAPIKey,
		from:   cfg.From,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// IsConfigured reports whether email delivery is enabled.
func (s *Service) IsConfigured() bool {
	return s != nil && s.apiKey != ""
}

// SendPasswordResetEmail sends the password reset link email.
func (s *Service) SendPasswordResetEmail(ctx context.Context, to, resetLink string) error {
	if !s.IsConfigured() {
		return fmt.Errorf("email service not configured")
	}

	text, htmlBody := passwordResetTemplate(resetLink, "1 hour")
	return s.send(ctx, to, "Reset your password", text, htmlBody)
}

func (s *Service) send(ctx context.Context, to, subject, textBody, htmlBody string) error {
	payload := map[string]any{
		"from":    s.from,
		"to":      to,
		"subject": subject,
		"text":    textBody,
		"html":    htmlBody,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendAPIURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Errorf("resend returned %s: %s", resp.Status, string(respBody))
}

func passwordResetTemplate(resetLink, expiresIn string) (string, string) {
	escapedLink := html.EscapeString(resetLink)
	escapedExpiresIn := html.EscapeString(expiresIn)

	textBody := fmt.Sprintf(
		"You requested a password reset. Click the link below to reset your password:\n\n%s\n\nThis link expires in %s. If you didn't request this, please ignore this email.",
		resetLink,
		expiresIn,
	)

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Reset Your Password</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
    .container { background: #f9f9f9; border-radius: 8px; padding: 30px; }
    h1 { color: #2563eb; font-size: 24px; margin-bottom: 20px; }
    .button { display: inline-block; background: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
    .link { word-break: break-all; color: #2563eb; }
    .footer { margin-top: 30px; font-size: 14px; color: #666; }
  </style>
</head>
<body>
  <div class="container">
    <h1>Reset Your Password</h1>
    <p>You requested a password reset. Click the button below to reset your password:</p>
    <a href="%s" class="button">Reset Password</a>
    <p>Or copy and paste this link into your browser:</p>
    <p class="link">%s</p>
    <div class="footer">
      <p>This link expires in <strong>%s</strong>.</p>
      <p>If you didn't request this, please ignore this email.</p>
    </div>
  </div>
</body>
</html>`, escapedLink, escapedLink, escapedExpiresIn)

	return textBody, htmlBody
}
