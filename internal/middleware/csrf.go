package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/irpanzy/Task-Forge/pkg/response"
)

const CSRFContextKey = "csrf"

func NewCSRF() fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		ContextKey:     CSRFContextKey,
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		CookieSecure:   false, // Set true jika sudah di environment HTTPS production
		CookieHTTPOnly: false, // False agar client JavaScript bisa membaca token
		Expiration:     1 * time.Hour,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return response.Error(c, fiber.StatusForbidden, "Invalid or missing CSRF token", err.Error())
		},
	})
}
