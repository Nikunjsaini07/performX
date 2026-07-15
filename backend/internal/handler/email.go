package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// senderEmail returns the verified Brevo sender email from the SENDER_EMAIL env var.
func senderEmail() string {
	if v := os.Getenv("SENDER_EMAIL"); v != "" {
		return v
	}
	return "no-reply@performx.app" // fallback
}

const brevoSenderName = "PerformX"

type brevoEmailAddress struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
}

type brevoEmailPayload struct {
	Sender      brevoEmailAddress   `json:"sender"`
	To          []brevoEmailAddress `json:"to"`
	Subject     string              `json:"subject"`
	HtmlContent string              `json:"htmlContent"`
}

// SendOTPEmail delivers a branded OTP email via Brevo's transactional REST API.
// purpose must be one of: "REGISTER", "LOGIN", "FORGOT_PASSWORD"
func SendOTPEmail(toEmail, toName, otpCode, purpose string) error {
	apiKey := os.Getenv("BREVO_API")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API key not set in environment")
	}

	subject, body := buildOTPEmailContent(toName, otpCode, purpose)

	payload := brevoEmailPayload{
		Sender: brevoEmailAddress{
			Name:  brevoSenderName,
			Email: senderEmail(),
		},
		To: []brevoEmailAddress{
			{Name: toName, Email: toEmail},
		},
		Subject:     subject,
		HtmlContent: body,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email via Brevo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("brevo API returned status %d", resp.StatusCode)
	}

	return nil
}

// buildOTPEmailContent returns the email subject and HTML body based on the OTP purpose.
func buildOTPEmailContent(name, otp, purpose string) (subject, html string) {
	if name == "" {
		name = "there"
	}

	var heading, actionLine, note string

	switch purpose {
	case "REGISTER":
		subject = "Verify your PerformX account"
		heading = "Welcome to PerformX! 🎉"
		actionLine = "Use the code below to verify your email and activate your account:"
		note = "This code expires in <strong>15 minutes</strong>. If you didn't create a PerformX account, you can safely ignore this email."

	case "LOGIN":
		subject = "Your PerformX login verification code"
		heading = "Login Verification"
		actionLine = "Use the code below to complete your sign-in:"
		note = "This code expires in <strong>15 minutes</strong>. If you didn't attempt to log in, please reset your password immediately."

	case "FORGOT_PASSWORD":
		subject = "Reset your PerformX password"
		heading = "Password Reset Request 🔐"
		actionLine = "Use the code below to reset your password:"
		note = "This code expires in <strong>15 minutes</strong>. If you didn't request a password reset, you can safely ignore this email."

	default:
		subject = "Your PerformX verification code"
		heading = "Verification Code"
		actionLine = "Use the code below to continue:"
		note = "This code expires in <strong>15 minutes</strong>."
	}

	html = fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>%s</title>
</head>
<body style="margin:0;padding:0;background-color:#0f0f13;font-family:'Segoe UI',Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#0f0f13;padding:40px 20px;">
    <tr>
      <td align="center">
        <table width="560" cellpadding="0" cellspacing="0" style="max-width:560px;width:100%%;">

          <!-- Logo / Brand -->
          <tr>
            <td align="center" style="padding-bottom:32px;">
              <span style="font-size:26px;font-weight:800;color:#ffffff;letter-spacing:-0.5px;">
                Perform<span style="color:#6c63ff;">X</span>
              </span>
            </td>
          </tr>

          <!-- Card -->
          <tr>
            <td style="background:linear-gradient(135deg,#1a1a2e 0%%,#16213e 100%%);border-radius:16px;border:1px solid #2a2a3d;padding:40px 36px;">

              <!-- Heading -->
              <p style="margin:0 0 8px;font-size:22px;font-weight:700;color:#ffffff;">%s</p>
              <p style="margin:0 0 28px;font-size:15px;color:#9fa3b1;">Hi %s,</p>
              <p style="margin:0 0 28px;font-size:15px;color:#c0c4d6;line-height:1.6;">%s</p>

              <!-- OTP Code Box -->
              <table width="100%%" cellpadding="0" cellspacing="0" style="margin-bottom:28px;">
                <tr>
                  <td align="center" style="background:linear-gradient(135deg,#6c63ff22,#a855f722);border:1px solid #6c63ff55;border-radius:12px;padding:24px;">
                    <span style="font-size:42px;font-weight:800;letter-spacing:12px;color:#ffffff;font-family:'Courier New',monospace;">%s</span>
                  </td>
                </tr>
              </table>

              <!-- Note -->
              <p style="margin:0;font-size:13px;color:#72758a;line-height:1.6;">%s</p>

            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td align="center" style="padding-top:24px;">
              <p style="margin:0;font-size:12px;color:#4a4d60;">
                &copy; 2025 PerformX. All rights reserved.<br/>
                <span style="color:#6c63ff;">performx.app</span>
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, subject, heading, name, actionLine, otp, note)

	return subject, html
}
