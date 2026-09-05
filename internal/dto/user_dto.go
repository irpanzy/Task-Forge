package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/model"
)

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	PublicID  uuid.UUID `json:"public_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token,omitempty"`
}

type PaginatedUsersResponse struct {
	Users       []UserResponse `json:"users"`
	TotalData   int64          `json:"total_data"`
	CurrentPage int            `json:"current_page"`
	TotalPages  int            `json:"total_pages"`
	Limit       int            `json:"limit"`
}

func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		PublicID:  u.PublicID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
