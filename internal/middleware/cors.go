package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/irpanzy/Task-Forge/internal/config"
)

func NewCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.CORSOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-CSRF-Token",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	})
}
