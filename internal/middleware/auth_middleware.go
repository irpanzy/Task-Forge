package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/irpanzy/Task-Forge/pkg/response"
	"github.com/irpanzy/Task-Forge/pkg/utils"
)

func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		tokenString = c.Cookies("access_token")

		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Authentication failed: token not found", nil)
		}

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Authentication failed: "+err.Error(), nil)
		}

		if publicID, ok := claims["public_id"].(string); ok {
			c.Locals("public_id", publicID)
		}
		if role, ok := claims["role"].(string); ok {
			c.Locals("role", role)
		}
		if email, ok := claims["email"].(string); ok {
			c.Locals("email", email)
		}
		if userID, ok := claims["user_id"]; ok {
			c.Locals("user_id", userID)
		}

		return c.Next()
	}
}

func RequireRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, _ := c.Locals("role").(string)
		if userRole == "" {
			return response.Error(c, fiber.StatusForbidden, "Access denied: undefined role", nil)
		}

		for _, role := range allowedRoles {
			if strings.EqualFold(userRole, role) {
				return c.Next()
			}
		}

		return response.Error(c, fiber.StatusForbidden, "Access denied: you do not have permission for this action", nil)
	}
}
