// Package repositories Description: Этот файл содержит репозиторий для пользователей.
package repositories

import (
	"merch-shop/internal/app/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=user_repository.go -destination=../../../mocks/repositories/user_repository.go

// UserRepositoryInterface описывает репозиторий для пользователей.
type UserRepositoryInterface interface {
	FindByUsername(tx *gorm.DB, username string) (*models.User, error)
	FindByID(tx *gorm.DB, userID uint) (*models.User, error)
	Save(tx *gorm.DB, user *models.User) error
}

// UserRepository репозиторий для пользователей.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый репозиторий для пользователей.
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// FindByUsername находит пользователя по имени.
func (r *UserRepository) FindByUsername(tx *gorm.DB, username string) (*models.User, error) {
	if tx == nil {
		tx = r.db
	}
	var user models.User
	result := tx.Where("username = ?", username).First(&user)
	return &user, result.Error
}

// FindByID находит пользователя по ID.
func (r *UserRepository) FindByID(tx *gorm.DB, userID uint) (*models.User, error) {
	if tx == nil {
		tx = r.db
	}
	var user models.User
	result := tx.Where("id = ?", userID).First(&user)
	return &user, result.Error
}

// Save сохраняет пользователя.
func (r *UserRepository) Save(tx *gorm.DB, user *models.User) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(user).Error
}
