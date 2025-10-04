package auth

import (
	"fmt"
	"varaden/server/config"
)

func (am *AuthModule) SendVerificationEmail(to, token string) error {
	subject := "Email Verification"

	body := fmt.Sprintf(`
Dear user,

To verify your email, enter this code: %s

If you did not create an account, then ignore this email.
	`, token)

	return am.email.SendEmail(to, subject, body)
}

func (am *AuthModule) SendResetPasswordEmail(to, token string) error {
	subject := "Reset password"

	resetPasswordURL := fmt.Sprintf("%s/reset-password?token=%s", config.FrontEndURL, token)
	body := fmt.Sprintf(`
Dear user,

To reset your password, click on this link: %s

If you did not request any password resets, then ignore this email.
`, resetPasswordURL)
	return am.email.SendEmail(to, subject, body)
}
