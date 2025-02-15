package repositories

import (
	"gorm.io/gorm"
	"merch-shop/internal/app/models"
)

type IUserRepository interface {
	FindByUsername(username string) (*models.User, error)
	Save(user *models.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	return &user, result.Error
}

func (r *UserRepository) Save(user *models.User) error {
	return r.db.Save(user).Error
}
