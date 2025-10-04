package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func FiberAppMiddlewares(app *fiber.App) {
	app.Use(LoggerConfig())
	app.Use(RecoverConfig())
}
