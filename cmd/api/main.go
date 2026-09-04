package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/irpanzy/Task-Forge/internal/config"
)

func main() {
	// 1. Load environment variables & database connection
	config.LoadEnv()
	config.ConnectDB()

	// 2. Initialize Fiber app
	app := fiber.New()

	// Base health check route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":     "TaskForge API",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// 3. Start server
	port := config.AppConfig.Port
	if port == "" {
		port = "3000"
	}

	log.Printf("Server TaskForge berjalan di http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
