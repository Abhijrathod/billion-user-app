package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/my-username/billion-user-app/pkg/config"
	"github.com/my-username/billion-user-app/pkg/database"
	"github.com/my-username/billion-user-app/pkg/jwtutils"
	"github.com/my-username/billion-user-app/pkg/logger"
	"github.com/my-username/billion-user-app/services/media-service/internal/domain"
	"github.com/my-username/billion-user-app/services/media-service/internal/handler"
	"github.com/my-username/billion-user-app/services/media-service/internal/repository"
	"github.com/my-username/billion-user-app/services/media-service/internal/service"
)

func main() {
	cfg, err := config.LoadConfig("../../.env")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	appLogger := logger.New("media-service")
	appLogger.Info().Msg("Starting media service")

	db, err := database.ConnectDB(cfg, "media_db")
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	if err := db.AutoMigrate(&domain.Media{}); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to migrate database")
	}

	jwtManager := jwtutils.NewJWTManager(cfg.JWTSecret, 15*time.Minute)

	mediaRepo := repository.NewMediaRepository(db)
	mediaService := service.NewMediaService(mediaRepo)
	mediaHandler := handler.NewMediaHandler(mediaService)

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
		return c.JSON(fiber.Map{"status": "ok", "service": "media-service"})
	})

	api := app.Group("/api/v1")

	// Protected routes (all media operations require auth)
	protected := api.Group("/media", handler.JWTMiddleware(jwtManager))
	protected.Post("/", mediaHandler.CreateMedia)
	protected.Get("/", mediaHandler.GetMyMedia)
	protected.Get("/:id", mediaHandler.GetMedia)
	protected.Delete("/:id", mediaHandler.DeleteMedia)
	protected.Post("/presigned-url", mediaHandler.GetPresignedURL)

	port := cfg.Port
	if port == "" {
		port = "3005"
	}
	appLogger.Info().Str("port", port).Msg("Starting server")
	if err := app.Listen(":" + port); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to start server")
	}
}

