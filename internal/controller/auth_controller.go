package controller

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/config"
	"github.com/irpanzy/Task-Forge/internal/dto"
	"github.com/irpanzy/Task-Forge/internal/middleware"
	"github.com/irpanzy/Task-Forge/internal/service"
	"github.com/irpanzy/Task-Forge/pkg/response"
	"github.com/irpanzy/Task-Forge/pkg/utils"
)

type AuthController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	GetCSRFToken(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}

type authController struct {
	userService service.UserService
}

func NewAuthController(userService service.UserService) AuthController {
	return &authController{userService: userService}
}

func (ctrl *authController) GetCSRFToken(c *fiber.Ctx) error {
	token, _ := c.Locals(middleware.CSRFContextKey).(string)
	if token == "" {
		token = c.Cookies("csrf_")
	}
	if token == "" {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve CSRF token", nil)
	}

	return response.Success(c, fiber.StatusOK, "CSRF token retrieved successfully", fiber.Map{
		"csrf_token": token,
	})
}

func (ctrl *authController) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON format", err.Error())
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Name == "" {
		return response.Error(c, fiber.StatusBadRequest, "Name is required", nil)
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return response.Error(c, fiber.StatusBadRequest, "Invalid email format", nil)
	}
	if len(req.Password) < 6 {
		return response.Error(c, fiber.StatusBadRequest, "Password must be at least 6 characters", nil)
	}

	res, err := ctrl.userService.Register(&req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Registration successful", res)
}

func (ctrl *authController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON format", err.Error())
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return response.Error(c, fiber.StatusBadRequest, "Email and password are required", nil)
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

	return response.Success(c, fiber.StatusOK, "Login successful", loginRes.User)
}

func (ctrl *authController) Logout(c *fiber.Ctx) error {
	c.ClearCookie("access_token", "refresh_token")
	return response.Success(c, fiber.StatusOK, "Logout successful", nil)
}

func (ctrl *authController) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Refresh token not found in cookie", nil)
	}

	claims, err := utils.VerifyToken(refreshToken)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired refresh token", nil)
	}

	publicIDStr, ok := claims["public_id"].(string)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid token payload", nil)
	}

	publicID, err := uuid.Parse(publicIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid public_id format", nil)
	}

	userRes, err := ctrl.userService.GetDetail(publicID)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "User not found", nil)
	}

	var userID int64
	if idVal, ok := claims["user_id"].(float64); ok {
		userID = int64(idVal)
	}

	newAccessToken, err := utils.GenerateToken(userID, userRes.Role, userRes.Email, userRes.PublicID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to generate new access token", nil)
	}

	accessDuration := utils.ParseDuration(config.AppConfig.JWTExpired, 15*time.Minute)
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Expires:  time.Now().Add(accessDuration),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return response.Success(c, fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"access_token": newAccessToken,
	})
}
