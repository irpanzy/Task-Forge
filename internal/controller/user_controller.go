package controller

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/dto"
	"github.com/irpanzy/Task-Forge/internal/service"
	"github.com/irpanzy/Task-Forge/pkg/response"
)

type UserController interface {
	GetProfile(c *fiber.Ctx) error
	GetUsers(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

// GetProfile returns the profile of the currently authenticated user
func (ctrl *userController) GetProfile(c *fiber.Ctx) error {
	publicIDStr, _ := c.Locals("public_id").(string)
	publicID, err := uuid.Parse(publicIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid user identity", nil)
	}

	user, err := ctrl.userService.GetDetail(publicID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Profile retrieved successfully", user)
}

// GetUsers retrieves all users with pagination and search filter - Admin only
func (ctrl *userController) GetUsers(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	result, err := ctrl.userService.GetUsers(search, page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve users", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Users retrieved successfully", result)
}

// GetUserByID retrieves user details by public_id
func (ctrl *userController) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID format", nil)
	}

	user, err := ctrl.userService.GetDetail(publicID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "User retrieved successfully", user)
}

// UpdateUser updates user details (name, email)
func (ctrl *userController) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID format", nil)
	}

	// Authorization check: non-admin can only update their own profile
	currentUserID, _ := c.Locals("public_id").(string)
	currentUserRole, _ := c.Locals("role").(string)

	if !strings.EqualFold(currentUserRole, "admin") && currentUserID != publicID.String() {
		return response.Error(c, fiber.StatusForbidden, "Access denied: You do not have permission to modify another user's data", nil)
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON format", err.Error())
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	updatedUser, err := ctrl.userService.Update(publicID, &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "User updated successfully", updatedUser)
}

// DeleteUser deletes a user by public_id - Admin only
func (ctrl *userController) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID format", nil)
	}

	currentUserID, _ := c.Locals("public_id").(string)
	if currentUserID == publicID.String() {
		return response.Error(c, fiber.StatusBadRequest, "Cannot delete your own active account", nil)
	}

	if err := ctrl.userService.Delete(publicID); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}
