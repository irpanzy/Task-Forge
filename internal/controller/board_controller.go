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

type BoardController interface {
	CreateBoard(c *fiber.Ctx) error
	GetBoards(c *fiber.Ctx) error
	GetBoardByID(c *fiber.Ctx) error
	UpdateBoard(c *fiber.Ctx) error
	DeleteBoard(c *fiber.Ctx) error
}

type boardController struct {
	boardService service.BoardService
}

func NewBoardController(boardService service.BoardService) BoardController {
	return &boardController{boardService: boardService}
}

// CreateBoard creates a new board for the authenticated user
func (ctrl *boardController) CreateBoard(c *fiber.Ctx) error {
	userPublicIDStr, _ := c.Locals("public_id").(string)
	userPublicID, err := uuid.Parse(userPublicIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid user identity", nil)
	}

	var req dto.CreateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON format", err.Error())
	}

	res, err := ctrl.boardService.CreateBoard(userPublicID, &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Board created successfully", res)
}

// GetBoards retrieves boards belonging to the authenticated user with pagination and search
func (ctrl *boardController) GetBoards(c *fiber.Ctx) error {
	userPublicIDStr, _ := c.Locals("public_id").(string)
	userPublicID, err := uuid.Parse(userPublicIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid user identity", nil)
	}

	search := strings.TrimSpace(c.Query("search"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	res, err := ctrl.boardService.GetUserBoards(userPublicID, search, page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve boards", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Boards retrieved successfully", res)
}

// GetBoardByID retrieves a specific board by public UUID
func (ctrl *boardController) GetBoardByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	boardPublicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid board ID format", nil)
	}

	userPublicIDStr, _ := c.Locals("public_id").(string)
	userPublicID, _ := uuid.Parse(userPublicIDStr)
	userRole, _ := c.Locals("role").(string)

	board, err := ctrl.boardService.GetBoardDetail(boardPublicID, userPublicID, userRole)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
		}
		return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Board retrieved successfully", board)
}

// UpdateBoard updates board details
func (ctrl *boardController) UpdateBoard(c *fiber.Ctx) error {
	idParam := c.Params("id")
	boardPublicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid board ID format", nil)
	}

	userPublicIDStr, _ := c.Locals("public_id").(string)
	userPublicID, _ := uuid.Parse(userPublicIDStr)
	userRole, _ := c.Locals("role").(string)

	var req dto.UpdateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON format", err.Error())
	}

	updatedBoard, err := ctrl.boardService.UpdateBoard(boardPublicID, userPublicID, userRole, &req)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Board updated successfully", updatedBoard)
}

// DeleteBoard deletes a board by public UUID
func (ctrl *boardController) DeleteBoard(c *fiber.Ctx) error {
	idParam := c.Params("id")
	boardPublicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid board ID format", nil)
	}

	userPublicIDStr, _ := c.Locals("public_id").(string)
	userPublicID, _ := uuid.Parse(userPublicIDStr)
	userRole, _ := c.Locals("role").(string)

	if err := ctrl.boardService.DeleteBoard(boardPublicID, userPublicID, userRole); err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Board deleted successfully", nil)
}
