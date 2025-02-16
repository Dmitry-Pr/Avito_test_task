package services

import (
	"errors"
	"fmt"
	"merch-shop/internal/app/models"
	"merch-shop/internal/app/repositories"
	"merch-shop/internal/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserService interface {
	Authenticate(username, password string) (string, error)
}

type UserService struct {
	repo repositories.IUserRepository
}

func NewUserService(repo repositories.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Authenticate(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return "", fmt.Errorf("ошибка хеширования пароля: %w", err)
			}

			newUser := &models.User{Username: username, Password: string(hashedPassword)}
			if err := s.repo.Save(newUser); err != nil {
				return "", fmt.Errorf("не удалось создать пользователя: %w", err)
			}
			user = newUser
		} else {
			return "", err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("неверные данные пользователя")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("не удалось создать токен: %w", err)
	}
	return token, nil
}
