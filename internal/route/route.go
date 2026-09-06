package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irpanzy/Task-Forge/internal/controller"
)

func SetupRoutes(app *fiber.App, authCtrl controller.AuthController, userCtrl controller.UserController, boardCtrl controller.BoardController) {
	api := app.Group("/api")

	RegisterAuthRoutes(api, authCtrl)
	RegisterUserRoutes(api, userCtrl)
	RegisterBoardRoutes(api, boardCtrl)
}
