package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irpanzy/Task-Forge/internal/controller"
)

func RegisterAuthRoutes(router fiber.Router, authCtrl controller.AuthController) {
	auth := router.Group("/auth")
	auth.Get("/csrf", authCtrl.GetCSRFToken)
	auth.Post("/register", authCtrl.Register)
	auth.Post("/login", authCtrl.Login)
	auth.Post("/logout", authCtrl.Logout)
	auth.Post("/refresh", authCtrl.RefreshToken)
}
