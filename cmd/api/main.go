package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"varaden/server/config"
	"varaden/server/internal/database"
	"varaden/server/internal/middlewares"
	"varaden/server/internal/modules"
	"varaden/server/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.AppConfig()
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	db, err := database.InitDatabase(cfg.DB)
	if err != nil {
		panic(err)
	}
	defer database.CloseDatabase(db)

	middlewares.FiberAppMiddlewares(app)
	modules.Setup(app, db, cfg)

	// Start server and handle graceful shutdown
	serverErrors := make(chan error, 1)
	go startServer(app, cfg.PortAddress, serverErrors)
	handleGracefulShutdown(ctx, app, serverErrors)
}

func startServer(app *fiber.App, address string, errs chan<- error) {
	if err := app.Listen(address); err != nil {
		errs <- fmt.Errorf("error starting server: %w", err)
	}
}

func handleGracefulShutdown(ctx context.Context, app *fiber.App, serverErrors <-chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	case <-quit:
		log.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error during server shutdown: %v", err)
		}
	case <-ctx.Done():
		log.Info("Server exiting due to context cancellation")
	}

	log.Info("Server exited")
}
