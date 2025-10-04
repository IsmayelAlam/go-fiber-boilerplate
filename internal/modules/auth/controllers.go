package auth

import (
	"context"
	"time"
	"varaden/server/config"
	authServices "varaden/server/internal/modules/auth/services"
	userServices "varaden/server/internal/modules/user/services"
	"varaden/server/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (am *AuthModule) register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	req := new(registerData)

	// Parse and validate request
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := am.validate.Struct(req); err != nil {
		return err
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Create user
	newUser, err := am.user.CreateUser(ctx, userServices.CreateUserParams{
		Email:        req.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return utils.DuplicateEntryError(err, "email")
	}

	// Create email verification token
	otp := utils.GenerateRandomNumber()

	newToken, err := am.token.CreateToken(ctx, authServices.CreateTokenParams{
		UserID:    newUser.ID,
		Token:     otp,
		Type:      authServices.TokenTypeEmailVerify,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return err
	}

	// Send verification email
	if err := am.SendVerificationEmail(newUser.Email, newToken.Token); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": newUser,
	})
}

func (am *AuthModule) login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	req := new(registerData)

	// Parse and validate request
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := am.validate.Struct(req); err != nil {
		return err
	}

	// Authenticate user
	user, err := am.user.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}
	if !user.IsActive {
		return fiber.NewError(fiber.StatusUnauthorized, "User account is deactivated. Contact support.")
	}
	if !user.VerifiedEmail {
		// Get token
		getToken, err := am.token.GetTokenByUserId(ctx, user.ID)

		if err == nil && getToken.ExpiresAt.Before(time.Now()) {
			// Send verification email
			if err := am.SendVerificationEmail(user.Email, getToken.Token); err != nil {
				return err
			}
			return c.JSON(fiber.Map{
				"data": fiber.Map{
					"id":             user.ID,
					"verified_email": false,
				},
			})
		} else {
			am.token.DeleteToken(ctx, getToken.ID)
		}

		// Create email verification token
		otp := utils.GenerateRandomNumber()
		newToken, err := am.token.CreateToken(ctx, authServices.CreateTokenParams{
			UserID:    user.ID,
			Token:     otp,
			Type:      authServices.TokenTypeEmailVerify,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		if err != nil {
			return err
		}

		// Send verification email
		if err := am.SendVerificationEmail(user.Email, newToken.Token); err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"data": fiber.Map{
				"id":             user.ID,
				"verified_email": false,
			},
		})
	}

	if user.LockedUntil.Valid && user.LockedUntil.Time.After(time.Now()) {
		return fiber.NewError(fiber.StatusTooManyRequests, "Account locked due to multiple failed login attempts.")
	}

	// Check password
	if matched := utils.CheckPasswordHash(req.Password, user.PasswordHash); !matched {
		am.user.IncrementFailedLogin(ctx, user.ID)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Reset failed login attempts
	if err := am.user.ResetFailedLogin(ctx, user.ID); err != nil {
		return err
	}

	// Generate JWT tokens
	tokens, err := am.jwt.GenerateToken(user.ID)
	if err != nil {
		return err
	}

	// Set refresh token in HTTP-only cookie
	am.jwt.SetRefreshCookie(c, tokens.RefreshToken)

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"email":          user.Email,
			"name":           user.Name,
			"id":             user.ID,
			"verified_email": user.VerifiedEmail,
			"aToken":         tokens.Token,
		},
	})
}

func (am *AuthModule) refreshTokens(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	req := new(refreshTokensData)
	// Parse and validate request
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := am.validate.Struct(req); err != nil {
		return err
	}

	if req.Logout {
		am.jwt.GetExpiredRefreshCookie(c)
		return c.JSON(fiber.Map{
			"data": fiber.Map{
				"message": "Logged out successfully",
			},
		})
	}

	refreshTokensData := c.Cookies(config.JWTConfig.RefreshCookieName)
	userId, err := am.jwt.RefreshTokenValidate(refreshTokensData)
	if err != nil {
		return c.JSON(fiber.Map{})
	}

	// Convert string to UUID
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(fiber.Map{})
	}

	// Authenticate user
	user, err := am.user.GetUserById(ctx, userUUID)
	if err != nil || !user.VerifiedEmail || !user.IsActive {
		return c.JSON(fiber.Map{})
	}

	tokens, err := am.jwt.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(fiber.Map{})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"email":          user.Email,
			"name":           user.Name,
			"id":             user.ID,
			"verified_email": user.VerifiedEmail,
			"aToken":         tokens.Token,
		},
	})
}

func (am *AuthModule) forgotPassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	req := new(forgotPasswordData)
	// Parse and validate request
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := am.validate.Struct(req); err != nil {
		return err
	}

	// Get user by email
	user, err := am.user.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// To prevent email enumeration, return the same response
		return c.JSON(fiber.Map{
			"data": fiber.Map{"message": "If a user with that email exists, a password reset email has been sent"},
		})
	}

	// Get existing token and delete it
	existingToken, err := am.token.GetTokenByUserId(ctx, user.ID)
	if err == nil {
		am.token.DeleteToken(ctx, existingToken.ID)
	}

	// Create password reset token
	resetToken := utils.GenerateRandomString(32)

	_, err = am.token.CreateToken(ctx, authServices.CreateTokenParams{
		UserID:    user.ID,
		Token:     resetToken,
		Type:      authServices.TokenTypePasswordReset,
		ExpiresAt: time.Now().Add(12 * time.Hour),
	})
	if err != nil {
		return err
	}

	if err := am.SendResetPasswordEmail(user.Email, resetToken); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{"message": "If a user with that email exists, a password reset email has been sent"},
	})
}

func (am *AuthModule) resetPassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	req := new(resetPasswordData)
	// Parse and validate request
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := am.validate.Struct(req); err != nil {
		return err
	}

	// Check if passwords match
	if req.Password != req.ConfirmPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Passwords do not match")
	}

	// Get and validate token
	token, err := am.token.GetTokenByCode(ctx, req.Token)
	if err != nil || token.ExpiresAt.Before(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid or expired token")
	}
	am.token.DeleteToken(ctx, token.ID)

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Update user's password
	err = am.user.UpdatePassword(ctx, userServices.UpdatePasswordParams{
		PasswordHash: passwordHash,
		ID:           token.UserID,
	})
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"message": "Password reset successfully",
		},
	})
}

func (am *AuthModule) sendVerificationEmail(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	req := new(resendVerifyEmailData)
	// Parse and validate request
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := am.validate.Struct(req); err != nil {
		return err
	}
	// Get token
	getToken, err := am.token.GetTokenByUserId(ctx, req.UserID)

	if err == nil {
		am.token.DeleteToken(ctx, getToken.ID)
	}
	// Create email verification token
	otp := utils.GenerateRandomNumber()

	newToken, err := am.token.CreateToken(ctx, authServices.CreateTokenParams{
		UserID:    req.UserID,
		Token:     otp,
		Type:      authServices.TokenTypeEmailVerify,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return err
	}
	// Get user
	user, err := am.user.GetUserById(ctx, req.UserID)
	if err != nil {
		return err
	}
	// Send verification email
	if err := am.SendVerificationEmail(user.Email, newToken.Token); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "Verification email sent successfully",
	})
}

func (am *AuthModule) verifyEmail(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	req := new(verifyEmailData)

	// Parse and validate request
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := am.validate.Struct(req); err != nil {
		return err
	}

	// Get token
	getToken, err := am.token.GetToken(ctx, authServices.GetTokenParams{
		Type:   authServices.TokenTypeEmailVerify,
		Token:  req.OTP,
		UserID: req.UserID,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid OTP")
	}
	if getToken.ExpiresAt.Before(time.Now()) {
		am.token.DeleteToken(ctx, getToken.ID)
		return fiber.NewError(fiber.StatusBadRequest, "OTP has expired. Resend OTP Code.")
	}

	// Verify user's email
	if err := am.user.VerifyUserEmail(ctx, req.UserID); err != nil {
		return err
	}

	// Delete token
	if err := am.token.DeleteToken(ctx, getToken.ID); err != nil {
		return err
	}

	// Generate JWT tokens
	tokens, _ := am.jwt.GenerateToken(req.UserID)

	// Set refresh token in HTTP-only cookie
	am.jwt.SetRefreshCookie(c, tokens.RefreshToken)

	return c.JSON(fiber.Map{
		"data": fiber.Map{"message": "Email verified successfully"},
	})
}

func (am *AuthModule) googleLogin(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "googleLogin: not implemented",
	})
}

func (am *AuthModule) googleCallback(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "googleCallback: not implemented",
	})
}
