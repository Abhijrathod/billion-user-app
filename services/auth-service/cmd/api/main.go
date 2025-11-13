package main

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/my-username/billion-user-app/pkg/config"
	"github.com/my-username/billion-user-app/pkg/database"
	"github.com/my-username/billion-user-app/pkg/jwtutils"
	"github.com/my-username/billion-user-app/pkg/kafkaclient"
	app_logger "github.com/my-username/billion-user-app/pkg/logger"

	"github.com/my-username/billion-user-app/services/auth-service/internal/domain"
	"github.com/my-username/billion-user-app/services/auth-service/internal/handler"
	"github.com/my-username/billion-user-app/services/auth-service/internal/repository"
	"github.com/my-username/billion-user-app/services/auth-service/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("../../.env")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger (custom logger)
	appLogger := app_logger.New("auth-service")
	appLogger.Info().Msg("Starting auth service")

	// Connect to database
	db, err := database.ConnectDB(cfg, "auth_db")
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Auto-migrate
	if err := db.AutoMigrate(&domain.User{}, &domain.RefreshToken{}); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// Initialize Kafka client
	var kafkaClient *kafkaclient.Client
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	if len(brokers) > 0 && brokers[0] != "" {
		kafkaClient, err = kafkaclient.NewClient(brokers)
		if err != nil {
			appLogger.Warn().Err(err).Msg("Failed to connect to Kafka, continuing without events")
		}
	}

	// Initialize JWT manager
	jwtManager := jwtutils.NewJWTManager(cfg.JWTSecret, 15*time.Minute)

	// Initialize repository
	authRepo := repository.NewAuthRepository(db)

	// Initialize service
	authService := service.NewAuthService(authRepo, jwtManager, kafkaClient)

	// Initialize handler
	authHandler := handler.NewAuthHandler(authService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(fiberlogger.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "auth-service"})
	})

	// Public routes
	api := app.Group("/api/v1")
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)
	api.Post("/refresh", authHandler.RefreshToken)
	api.Post("/logout", authHandler.Logout)

	// Protected routes
	protected := api.Group("/auth", handler.JWTMiddleware(jwtManager))
	protected.Get("/profile", authHandler.GetProfile)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "3001"
	}
	appLogger.Info().Str("port", port).Msg("Starting server")
	if err := app.Listen(":" + port); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to start server")
	}
}
