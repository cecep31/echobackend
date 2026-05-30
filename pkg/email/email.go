package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html"
	"mime/multipart"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"echobackend/config"
	"echobackend/pkg/queue"
)

const taskTypePasswordReset = "email:password_reset"

// Service sends application emails through SMTP.
type Service struct {
	host     string
	port     int
	username string
	password string
	from     string
	timeout  time.Duration
	useTLS   bool
	queue    *queue.Service
}

type passwordResetPayload struct {
	To        string `json:"to"`
	ResetLink string `json:"reset_link"`
}

// NewService creates an SMTP-backed email service and registers its background tasks.
func NewService(cfg config.EmailConfig, taskQueue *queue.Service) *Service {
	service := &Service{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUsername,
		password: cfg.SMTPPassword,
		from:     cfg.From,
		timeout:  cfg.Timeout,
		useTLS:   cfg.UseTLS,
		queue:    taskQueue,
	}
	service.registerQueueHandlers()
	return service
}

func (s *Service) registerQueueHandlers() {
	if s == nil || s.queue == nil {
		return
	}
	s.queue.Handle(taskTypePasswordReset, s.handlePasswordResetTask)
}

// IsConfigured reports whether queued email delivery is enabled.
func (s *Service) IsConfigured() bool {
	return s != nil && s.hasSMTPConfig() && s.queue != nil && s.queue.IsConfigured()
}

func (s *Service) hasSMTPConfig() bool {
	return s != nil && s.host != "" && s.port > 0 && s.from != ""
}

// Close is kept for DI cleanup compatibility.
func (s *Service) Close() error {
	return nil
}

// EnqueuePasswordResetEmail queues a password reset email for Asynq delivery.
func (s *Service) EnqueuePasswordResetEmail(to, resetLink string) error {
	if !s.IsConfigured() {
		return fmt.Errorf("email service not configured")
	}

	payload := passwordResetPayload{To: to, ResetLink: resetLink}
	return s.queue.EnqueueJSON(taskTypePasswordReset, payload, queue.TaskOptions{Timeout: s.timeout})
}

func (s *Service) handlePasswordResetTask(ctx context.Context, payloadBytes []byte) error {
	var payload passwordResetPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w: %w", err, queue.SkipRetry)
	}

	if payload.To == "" || payload.ResetLink == "" {
		return fmt.Errorf("invalid password reset payload: %w", queue.SkipRetry)
	}

	return s.SendPasswordResetEmail(ctx, payload.To, payload.ResetLink)
}

// SendPasswordResetEmail sends the password reset link email.
func (s *Service) SendPasswordResetEmail(ctx context.Context, to, resetLink string) error {
	if !s.hasSMTPConfig() {
		return fmt.Errorf("email service not configured")
	}

	text, htmlBody := passwordResetTemplate(resetLink, "1 hour")
	return s.send(ctx, to, "Reset your password", text, htmlBody)
}

func (s *Service) send(ctx context.Context, to, subject, textBody, htmlBody string) error {
	message, err := buildMessage(s.from, to, subject, textBody, htmlBody)
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", s.host, s.port)
	dialer := net.Dialer{Timeout: s.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	if s.timeout > 0 {
		if err := conn.SetDeadline(time.Now().Add(s.timeout)); err != nil {
			return err
		}
	}

	var client *smtp.Client
	if s.useTLS {
		tlsConn := tls.Client(conn, &tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12})
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			return err
		}
		client, err = smtp.NewClient(tlsConn, s.host)
	} else {
		client, err = smtp.NewClient(conn, s.host)
	}
	if err != nil {
		return err
	}
	defer client.Close()

	if !s.useTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12}
			if err := client.StartTLS(tlsConfig); err != nil {
				return err
			}
		}
	}

	if s.username != "" || s.password != "" {
		auth := smtp.PlainAuth("", s.username, s.password, s.host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	fromAddress, err := parseAddress(s.from)
	if err != nil {
		return err
	}
	toAddress, err := parseAddress(to)
	if err != nil {
		return err
	}

	if err := client.Mail(fromAddress); err != nil {
		return err
	}
	if err := client.Rcpt(toAddress); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func buildMessage(from, to, subject, textBody, htmlBody string) ([]byte, error) {
	fromAddress, err := mail.ParseAddress(from)
	if err != nil {
		return nil, err
	}
	toAddress, err := mail.ParseAddress(to)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)

	textPart, err := multipartWriter.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/plain; charset=utf-8"},
		"Content-Transfer-Encoding": {"8bit"},
	})
	if err != nil {
		return nil, err
	}
	if _, err := textPart.Write([]byte(textBody)); err != nil {
		return nil, err
	}

	htmlPart, err := multipartWriter.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/html; charset=utf-8"},
		"Content-Transfer-Encoding": {"8bit"},
	})
	if err != nil {
		return nil, err
	}
	if _, err := htmlPart.Write([]byte(htmlBody)); err != nil {
		return nil, err
	}

	if err := multipartWriter.Close(); err != nil {
		return nil, err
	}

	var message bytes.Buffer
	headers := [][2]string{
		{"From", fromAddress.String()},
		{"To", toAddress.String()},
		{"Subject", sanitizeHeader(subject)},
		{"MIME-Version", "1.0"},
		{"Content-Type", fmt.Sprintf("multipart/alternative; boundary=%q", multipartWriter.Boundary())},
	}
	for _, header := range headers {
		message.WriteString(header[0])
		message.WriteString(": ")
		message.WriteString(header[1])
		message.WriteString("\r\n")
	}
	message.WriteString("\r\n")
	message.Write(body.Bytes())

	return message.Bytes(), nil
}

func parseAddress(raw string) (string, error) {
	address, err := mail.ParseAddress(raw)
	if err != nil {
		return "", err
	}
	return address.Address, nil
}

func sanitizeHeader(value string) string {
	value = strings.ReplaceAll(value, "\r", "")
	return strings.ReplaceAll(value, "\n", "")
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
