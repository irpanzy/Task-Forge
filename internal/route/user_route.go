package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irpanzy/Task-Forge/internal/controller"
	"github.com/irpanzy/Task-Forge/internal/middleware"
)

func RegisterUserRoutes(router fiber.Router, userCtrl controller.UserController) {
	users := router.Group("/users", middleware.Authenticate())
	users.Get("/me", userCtrl.GetProfile)
	users.Get("/", middleware.RequireRoles("admin"), userCtrl.GetUsers)
	users.Get("/:id", userCtrl.GetUserByID)
	users.Put("/:id", userCtrl.UpdateUser)
	users.Delete("/:id", middleware.RequireRoles("admin"), userCtrl.DeleteUser)
}
