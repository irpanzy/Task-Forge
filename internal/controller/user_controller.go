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

// GetProfile mengembalikan detail profile pengguna yang sedang login
func (ctrl *userController) GetProfile(c *fiber.Ctx) error {
	publicIDStr, _ := c.Locals("public_id").(string)
	publicID, err := uuid.Parse(publicIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Identitas pengguna tidak valid", nil)
	}

	user, err := ctrl.userService.GetDetail(publicID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mendapatkan data profil", user)
}

// GetUsers mengambil daftar seluruh pengguna (pagination) - Admin only
func (ctrl *userController) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	result, err := ctrl.userService.GetUsers(page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengambil data pengguna", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil data pengguna", result)
}

// GetUserByID mengambil detail user berdasarkan public_id
func (ctrl *userController) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format ID user tidak valid", nil)
	}

	user, err := ctrl.userService.GetDetail(publicID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mendapatkan data user", user)
}

// UpdateUser memperbarui data user (nama, email)
func (ctrl *userController) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format ID user tidak valid", nil)
	}

	// Cek otorisasi: non-admin hanya boleh mengubah profil miliknya sendiri
	currentUserID, _ := c.Locals("public_id").(string)
	currentUserRole, _ := c.Locals("role").(string)

	if !strings.EqualFold(currentUserRole, "admin") && currentUserID != publicID.String() {
		return response.Error(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak memiliki izin untuk mengubah data user lain", nil)
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format JSON tidak valid", err.Error())
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	updatedUser, err := ctrl.userService.Update(publicID, &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Berhasil memperbarui data user", updatedUser)
}

// DeleteUser menghapus user berdasarkan public_id - Admin only
func (ctrl *userController) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format ID user tidak valid", nil)
	}

	currentUserID, _ := c.Locals("public_id").(string)
	if currentUserID == publicID.String() {
		return response.Error(c, fiber.StatusBadRequest, "Tidak dapat menghapus akun Anda sendiri yang sedang aktif", nil)
	}

	if err := ctrl.userService.Delete(publicID); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "User berhasil dihapus", nil)
}
