package main

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/my-username/billion-user-app/pkg/config"
	"github.com/my-username/billion-user-app/pkg/database"
	"github.com/my-username/billion-user-app/pkg/jwtutils"
	"github.com/my-username/billion-user-app/pkg/kafkaclient"
	"github.com/my-username/billion-user-app/pkg/logger"
	"github.com/my-username/billion-user-app/services/user-service/internal/domain"
	"github.com/my-username/billion-user-app/services/user-service/internal/handler"
	"github.com/my-username/billion-user-app/services/user-service/internal/repository"
	"github.com/my-username/billion-user-app/services/user-service/internal/service"
)

func main() {
	cfg, err := config.LoadConfig("../../.env")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	appLogger := logger.New("user-service")
	appLogger.Info().Msg("Starting user service")

	db, err := database.ConnectDB(cfg, "user_db")
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to migrate database")
	}

	var kafkaClient *kafkaclient.Client
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	if len(brokers) > 0 && brokers[0] != "" {
		kafkaClient, err = kafkaclient.NewClient(brokers)
		if err != nil {
			appLogger.Warn().Err(err).Msg("Failed to connect to Kafka, continuing without events")
		}
	}

	jwtManager := jwtutils.NewJWTManager(cfg.JWTSecret, 15*time.Minute)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, kafkaClient)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "user-service"})
	})

	api := app.Group("/api/v1")

	// Public routes
	api.Get("/users/:id", userHandler.GetUser)
	api.Get("/users/username/:username", userHandler.GetUserByUsername)
	api.Get("/users", userHandler.ListUsers)
	api.Get("/users/search", userHandler.SearchUsers)

	// Protected routes
	protected := api.Group("/users", handler.JWTMiddleware(jwtManager))
	protected.Post("/", userHandler.CreateUser)
	protected.Put("/:id", userHandler.UpdateUser)
	protected.Delete("/:id", userHandler.DeleteUser)

	port := cfg.Port
	if port == "" {
		port = "3002"
	}
	appLogger.Info().Str("port", port).Msg("Starting server")
	if err := app.Listen(":" + port); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to start server")
	}
}
