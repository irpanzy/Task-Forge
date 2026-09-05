package controller

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/irpanzy/Task-Forge/internal/config"
	"github.com/irpanzy/Task-Forge/internal/dto"
	"github.com/irpanzy/Task-Forge/internal/service"
	"github.com/irpanzy/Task-Forge/pkg/response"
	"github.com/irpanzy/Task-Forge/pkg/utils"
)

type AuthController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	GetCSRFToken(c *fiber.Ctx) error
}

type authController struct {
	userService service.UserService
}

func NewAuthController(userService service.UserService) AuthController {
	return &authController{userService: userService}
}

func (ctrl *authController) GetCSRFToken(c *fiber.Ctx) error {
	token, ok := c.Locals(csrf.ConfigDefault.ContextKey).(string)
	if !ok || token == "" {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mendapatkan CSRF token", nil)
	}

	return response.Success(c, fiber.StatusOK, "CSRF token berhasil diambil", fiber.Map{
		"csrf_token": token,
	})
}

func (ctrl *authController) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format JSON tidak valid", err.Error())
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Name == "" {
		return response.Error(c, fiber.StatusBadRequest, "Nama wajib diisi", nil)
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return response.Error(c, fiber.StatusBadRequest, "Format email tidak valid", nil)
	}
	if len(req.Password) < 6 {
		return response.Error(c, fiber.StatusBadRequest, "Password minimal 6 karakter", nil)
	}

	res, err := ctrl.userService.Register(&req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Registrasi berhasil", res)
}

func (ctrl *authController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format JSON tidak valid", err.Error())
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return response.Error(c, fiber.StatusBadRequest, "Email dan password wajib diisi", nil)
	}

	loginRes, err := ctrl.userService.Login(&req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	accessDuration := utils.ParseDuration(config.AppConfig.JWTExpired, 15*time.Minute)
	refreshDuration := utils.ParseDuration(config.AppConfig.RefreshTokenExpired, 7*24*time.Hour)

	// 1. Simpan Access Token ke Cookie HTTPOnly
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    loginRes.AccessToken,
		Expires:  time.Now().Add(accessDuration),
		HTTPOnly: true,
		Secure:   false, // Set true jika sudah di environment HTTPS production
		SameSite: "Lax",
		Path:     "/",
	})

	// 2. Simpan Refresh Token ke Cookie HTTPOnly
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    loginRes.RefreshToken,
		Expires:  time.Now().Add(refreshDuration),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return response.Success(c, fiber.StatusOK, "Login berhasil", loginRes.User)
}

func (ctrl *authController) Logout(c *fiber.Ctx) error {
	c.ClearCookie("access_token", "refresh_token")
	return response.Success(c, fiber.StatusOK, "Logout berhasil", nil)
}
