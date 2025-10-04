package main

import (
	"log"
	"varaden/server/internal/middlewares"
	_ "varaden/server/temp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title						Varaden API
// @version					1.0
// @description				The number #1 renting platform in Bangladesh for renting out your properties and finding your next home.
// @termsOfService				http://dev.varaden.io/terms/
// @contact.name				API Support
// @contact.email				apiHelp@varaden.com
// @license.name				Proprietary
// @host						localhost:8080
// @BasePath					/api/v1
// @securityDefinitions.apikey	JWT
// @in							header
// @name						Authorization
func main() {
	app := fiber.New()

	middlewares.FiberAppMiddlewares(app)
	app.Get("*", swagger.HandlerDefault)

	if err := app.Listen(":1227"); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
