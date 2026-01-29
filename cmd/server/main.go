package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alifdwt/techtest-unsia-be/internal/config"
	"github.com/alifdwt/techtest-unsia-be/internal/db"
	"github.com/alifdwt/techtest-unsia-be/internal/transport"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// @title Technical Test Unsia API
// @version 1.0
// @description API for managing LMS
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the token.

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	app := fiber.New(fiber.Config{
		AppName: "Technical Test UNSIA Backend",
	})

	transport.RegisterRoutes(app, queries)

	// ==== START SERVER ====
	go func() {
		if err := app.Listen(cfg.Server.Port); err != nil {
			log.Printf("fiber stopped: %v", err)
		}
	}()
	log.Printf("server running on %s", cfg.Server.Port)

	// ==== SIGNAL HANDLING ====
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,    // SIGINT (Ctrl+C)
		syscall.SIGTERM, // Docker / K8s / VSCode
	)

	<-quit
	log.Println("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("server shut down gracefully")
}
