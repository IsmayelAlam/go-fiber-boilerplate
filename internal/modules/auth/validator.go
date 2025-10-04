package auth

import "github.com/google/uuid"

type registerData struct {
	Email    string `json:"email" validate:"required,email,max=250" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8,max=100,password" example:"password1"`
}

type forgotPasswordData struct {
	Email string `json:"email" validate:"required,email,max=250" example:"user@example.com"`
}

type resetPasswordData struct {
	Password        string `json:"new_password" validate:"required,min=8,max=100,password" example:"password1"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=100,password" example:"password1"`
	Token           string `json:"token" validate:"required,len=32" example:"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"`
}

type resendVerifyEmailData struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid,max=250" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type verifyEmailData struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid,max=250" example:"550e8400-e29b-41d4-a716-446655440000"`
	OTP    string    `json:"otp" validate:"required,len=6" example:"123456"`
}

type refreshTokensData struct {
	Logout bool `json:"logout" example:"true"`
}
