package user

import (
	"database/sql"
	"varaden/server/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserModule struct {
	db       *sql.DB
	route    fiber.Router
	validate *validator.Validate
}

func RegisterUserModule(route fiber.Router, db *sql.DB) *UserModule {
	return &UserModule{
		db:       db,
		route:    route,
		validate: utils.Validator(),
	}
}
