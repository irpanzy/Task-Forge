package repository

import (
	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByPublicID(publicID uuid.UUID) (*model.User, error)
	FindByInternalID(internalID int64) (*model.User, error)
	FindAll(search string, offset, limit int) ([]model.User, int64, error)
	Update(user *model.User) error
	Delete(publicID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPublicID(publicID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Where("public_id = ?", publicID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByInternalID(internalID int64) (*model.User, error) {
	var user model.User
	err := r.db.Where("internal_id = ?", internalID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(search string, offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(publicID uuid.UUID) error {
	return r.db.Where("public_id = ?", publicID).Delete(&model.User{}).Error
}
