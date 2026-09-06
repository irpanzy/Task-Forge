package service

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/dto"
	"github.com/irpanzy/Task-Forge/internal/model"
	"github.com/irpanzy/Task-Forge/internal/repository"
	"github.com/irpanzy/Task-Forge/pkg/utils"
	"gorm.io/gorm"
)

type UserService interface {
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetDetail(publicID uuid.UUID) (*dto.UserResponse, error)
	GetUsers(search string, page, limit int) (*dto.PaginatedUsersResponse, error)
	Update(publicID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(publicID uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email is already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	newUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     model.RoleUser,
	}

	if err := s.userRepo.Create(&newUser); err != nil {
		return nil, err
	}

	res := dto.ToUserResponse(&newUser)
	return &res, nil
}

func (s *userService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateToken(user.InternalID, string(user.Role), user.Email, user.PublicID)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.InternalID, user.PublicID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &dto.LoginResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) GetDetail(publicID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByPublicID(publicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	res := dto.ToUserResponse(user)
	return &res, nil
}

func (s *userService) GetUsers(search string, page, limit int) (*dto.PaginatedUsersResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, totalData, err := s.userRepo.FindAll(search, offset, limit)
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponse
	for _, u := range users {
		userResponses = append(userResponses, dto.ToUserResponse(&u))
	}

	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	return &dto.PaginatedUsersResponse{
		Users:       userResponses,
		TotalData:   totalData,
		CurrentPage: page,
		TotalPages:  totalPages,
		Limit:       limit,
	}, nil
}

func (s *userService) Update(publicID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByPublicID(publicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		existing, err := s.userRepo.FindByEmail(req.Email)
		if err == nil && existing != nil {
			return nil, errors.New("new email is already in use by another account")
		}
		user.Email = req.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	res := dto.ToUserResponse(user)
	return &res, nil
}

func (s *userService) Delete(publicID uuid.UUID) error {
	_, err := s.userRepo.FindByPublicID(publicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.Delete(publicID)
}
