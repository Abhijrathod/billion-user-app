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
	"github.com/my-username/billion-user-app/services/product-service/internal/domain"
	"github.com/my-username/billion-user-app/services/product-service/internal/handler"
	"github.com/my-username/billion-user-app/services/product-service/internal/repository"
	"github.com/my-username/billion-user-app/services/product-service/internal/service"
)

func main() {
	cfg, err := config.LoadConfig("../../.env")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	appLogger := logger.New("product-service")
	appLogger.Info().Msg("Starting product service")

	db, err := database.ConnectDB(cfg, "product_db")
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	if err := db.AutoMigrate(&domain.Product{}); err != nil {
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

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo, kafkaClient)
	productHandler := handler.NewProductHandler(productService)

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
		return c.JSON(fiber.Map{"status": "ok", "service": "product-service"})
	})

	api := app.Group("/api/v1")

	// Public routes
	api.Get("/products/:id", productHandler.GetProduct)
	api.Get("/products", productHandler.ListProducts)
	api.Get("/products/search", productHandler.SearchProducts)
	api.Get("/products/category/:category", productHandler.GetProductsByCategory)

	// Protected routes
	protected := api.Group("/products", handler.JWTMiddleware(jwtManager))
	protected.Post("/", productHandler.CreateProduct)
	protected.Put("/:id", productHandler.UpdateProduct)
	protected.Delete("/:id", productHandler.DeleteProduct)

	port := cfg.Port
	if port == "" {
		port = "3003"
	}
	appLogger.Info().Str("port", port).Msg("Starting server")
	if err := app.Listen(":" + port); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to start server")
	}
}

