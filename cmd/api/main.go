package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DMaryanskiy/go-idk/internal/config"
	"github.com/DMaryanskiy/go-idk/internal/handler"
	"github.com/DMaryanskiy/go-idk/internal/middleware"
	"github.com/DMaryanskiy/go-idk/internal/repository"
	"github.com/DMaryanskiy/go-idk/internal/service"
	"github.com/DMaryanskiy/go-idk/internal/validator"
	"github.com/DMaryanskiy/go-idk/pkg/database"
	"github.com/DMaryanskiy/go-idk/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"go.uber.org/zap"
)

func main() {
	// Init logger
	log := logger.New()
	defer func () {
		err := log.Sync()
		if err != nil {
			log.Fatal("Failed to flush logger", zap.Error(err))
		}
	}()

	// Load configuration
	cfg := config.Load()

	// Init database
	db, err := database.New(
		cfg.DatabaseURL,
		cfg.DBMaxOpenConns,
		cfg.DBMaxIdleConns,
		cfg.DBConnMaxLifetime,
	)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			log.Fatal("Failed to close db", zap.Error(errDB))
		}
	}()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Init layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, log)
	val := validator.New()
	userHandler := handler.NewUserHandler(userService, val, log)

	// Init Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler(log),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	})

	// Middleware
	app.Use(requestid.New())
	app.Use(middleware.Logger(log))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.CORSOrigins},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.RateLimitMax,
		Expiration: cfg.RateLimitExpiration,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate Limiter Exceeded",
			})
		},
	}))

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// API Routes
	api := app.Group("/api/v1")
	userHandler.RegisterRoutes(api)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatal("Failed tp start server", zap.Error(err))
		}
	}()

	log.Info("Server started", zap.String("port", cfg.Port))

	<-quit
	log.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}

func customErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		logger.Error("Request error",
			zap.String("request_id", c.Locals("requestid").(string)),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Int("status", code),
			zap.Error(err),
		)

		return c.Status(code).JSON(fiber.Map{
			"error":      message,
			"request_id": c.Locals("requestid"),
		})
	}
}
