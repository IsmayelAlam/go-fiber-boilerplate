package modules

import (
	"database/sql"
	"varaden/server/config"
	"varaden/server/internal/modules/auth"
	healthCheck "varaden/server/internal/modules/health_check"
	"varaden/server/internal/modules/user"
	"varaden/server/internal/services"
	"varaden/server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, db *sql.DB, config config.AllConfig) {
	v1Group := app.Group("/api/v1")
	emailService := services.NewEmailService(&config.SMTP)

	user.RegisterUserModule(v1Group, db).SetupRoutes()
	auth.RegisterAuthModule(v1Group, db, emailService).SetupRoutes()
	healthCheck.RegisterHealthCheckModule(v1Group, db).SetupRoutes()

	// 404 Handler
	app.Use(utils.NotFoundHandler)
}
