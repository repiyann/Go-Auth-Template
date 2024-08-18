package main

import (
	"log"
	"os"
	database "template-auth/app/config"
	"template-auth/app/routes"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize .env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect database
	database.ConnectDB()

	// Initialize fiber
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,   // This is for making fiber more faster
		JSONDecoder: json.Unmarshal, // Based on their documentation
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Generic error handler
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fiber.Map{
					"errors": "Internal Server Error!",
				},
			})
		},
	})

	// Middlewares
	app.Use(etag.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	// Initialize all routes
	routes.SetupRoutes(app)

	// 404 Handlers
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fiber.Map{
				"errors": "The endpoint you are looking for does not exist. Please check the API documentation.",
			},
		})
	})

	// Initialize server
	PORT := os.Getenv("PORT")
	log.Printf("Starting server on port %s", PORT)
	log.Fatal(app.Listen(":" + PORT))
}
