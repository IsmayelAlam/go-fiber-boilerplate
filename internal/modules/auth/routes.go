package auth

func (am *AuthModule) SetupRoutes() {
	auth := am.route.Group("/auth")

	auth.Post("/register", am.register)
	auth.Post("/login", am.login)
	auth.Post("/refresh", am.refreshTokens)
	auth.Post("/forgot-password", am.forgotPassword)
	auth.Post("/reset-password", am.resetPassword)
	auth.Post("/send-verification-email", am.sendVerificationEmail)
	auth.Post("/verify-email", am.verifyEmail)
	auth.Get("/google", am.googleLogin)
	auth.Get("/google-callback", am.googleCallback)
}
