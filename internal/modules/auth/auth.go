package auth

import (
	"database/sql"
	"varaden/server/config"
	authServices "varaden/server/internal/modules/auth/services"
	userServices "varaden/server/internal/modules/user/services"
	"varaden/server/internal/services"
	"varaden/server/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthModule struct {
	db       *sql.DB
	route    fiber.Router
	validate *validator.Validate
	email    services.EmailService
	token    *authServices.Queries
	user     *userServices.Queries
	jwt      *utils.JWTConfig
}

func RegisterAuthModule(route fiber.Router, db *sql.DB, emailService services.EmailService) *AuthModule {
	jwtConfig := config.JWTConfig

	return &AuthModule{
		db:       db,
		route:    route,
		email:    emailService,
		validate: utils.Validator(),
		token:    authServices.New(db),
		user:     userServices.New(db),
		jwt:      jwtConfig,
	}
}
