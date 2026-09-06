package service

import (
	"errors"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/dto"
	"github.com/irpanzy/Task-Forge/internal/model"
	"github.com/irpanzy/Task-Forge/internal/repository"
	"gorm.io/gorm"
)

type BoardService interface {
	CreateBoard(ownerPublicID uuid.UUID, req *dto.CreateBoardRequest) (*dto.BoardResponse, error)
	GetUserBoards(ownerPublicID uuid.UUID, search string, page, limit int) (*dto.PaginatedBoardsResponse, error)
	GetBoardDetail(boardPublicID, userPublicID uuid.UUID, userRole string) (*dto.BoardResponse, error)
	UpdateBoard(boardPublicID, userPublicID uuid.UUID, userRole string, req *dto.UpdateBoardRequest) (*dto.BoardResponse, error)
	DeleteBoard(boardPublicID, userPublicID uuid.UUID, userRole string) error
}

type boardService struct {
	boardRepo repository.BoardRepository
	userRepo  repository.UserRepository
}

func NewBoardService(boardRepo repository.BoardRepository, userRepo repository.UserRepository) BoardService {
	return &boardService{
		boardRepo: boardRepo,
		userRepo:  userRepo,
	}
}

func (s *boardService) CreateBoard(ownerPublicID uuid.UUID, req *dto.CreateBoardRequest) (*dto.BoardResponse, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return nil, errors.New("board title is required")
	}

	user, err := s.userRepo.FindByPublicID(ownerPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("owner user not found")
		}
		return nil, err
	}

	newBoard := model.Board{
		OwnerID:       user.InternalID,
		OwnerPublicID: ownerPublicID,
		Title:         req.Title,
		Description:   strings.TrimSpace(req.Description),
		DueDate:       req.DueDate,
	}

	if err := s.boardRepo.Create(&newBoard); err != nil {
		return nil, err
	}

	res := dto.ToBoardResponse(&newBoard)
	return &res, nil
}

func (s *boardService) GetUserBoards(ownerPublicID uuid.UUID, search string, page, limit int) (*dto.PaginatedBoardsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	boards, totalData, err := s.boardRepo.FindAllByOwner(ownerPublicID, search, offset, limit)
	if err != nil {
		return nil, err
	}

	var boardResponses []dto.BoardResponse
	for _, b := range boards {
		boardResponses = append(boardResponses, dto.ToBoardResponse(&b))
	}

	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	return &dto.PaginatedBoardsResponse{
		Boards:      boardResponses,
		TotalData:   totalData,
		CurrentPage: page,
		TotalPages:  totalPages,
		Limit:       limit,
	}, nil
}

func (s *boardService) GetBoardDetail(boardPublicID, userPublicID uuid.UUID, userRole string) (*dto.BoardResponse, error) {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, err
	}

	// Authorization check: non-admin can only access boards they own
	if !strings.EqualFold(userRole, "admin") && board.OwnerPublicID != userPublicID {
		return nil, errors.New("access denied: you do not have permission to view this board")
	}

	res := dto.ToBoardResponse(board)
	return &res, nil
}

func (s *boardService) UpdateBoard(boardPublicID, userPublicID uuid.UUID, userRole string, req *dto.UpdateBoardRequest) (*dto.BoardResponse, error) {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, err
	}

	// Authorization check: non-admin can only update boards they own
	if !strings.EqualFold(userRole, "admin") && board.OwnerPublicID != userPublicID {
		return nil, errors.New("access denied: you do not have permission to modify this board")
	}

	if req.Title != "" {
		board.Title = strings.TrimSpace(req.Title)
	}
	if req.Description != "" {
		board.Description = strings.TrimSpace(req.Description)
	}
	if req.DueDate != nil {
		board.DueDate = req.DueDate
	}

	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}

	res := dto.ToBoardResponse(board)
	return &res, nil
}

func (s *boardService) DeleteBoard(boardPublicID, userPublicID uuid.UUID, userRole string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found")
		}
		return err
	}

	// Authorization check: non-admin can only delete boards they own
	if !strings.EqualFold(userRole, "admin") && board.OwnerPublicID != userPublicID {
		return errors.New("access denied: you do not have permission to delete this board")
	}

	return s.boardRepo.Delete(boardPublicID)
}
