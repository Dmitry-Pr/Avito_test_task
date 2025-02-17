// Package services Description: Описание сервиса для работы с пользователями.
package services

import (
	"errors"
	"fmt"
	"merch-shop/internal/app/models"
	"merch-shop/internal/app/repositories"
	"merch-shop/internal/pkg/utils"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//go:generate mockgen -source=user_service.go -destination=../../../mocks/services/user_service.go

// UserServiceInterface описывает сервис для работы с пользователями.
type UserServiceInterface interface {
	Authenticate(username, password string) (string, error)
}

// UserService сервис для работы с пользователями.
type UserService struct {
	repo repositories.UserRepositoryInterface
}

// NewUserService создает новый сервис для работы с пользователями.
func NewUserService(repo repositories.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

// Authenticate аутентифицирует пользователя.
func (s *UserService) Authenticate(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(nil, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return "", fmt.Errorf("ошибка хеширования пароля: %w", err)
			}

			newUser := &models.User{Username: username, Password: string(hashedPassword)}
			if err := s.repo.Save(nil, newUser); err != nil {
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

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY переменная среды не найдена")
	}
	token, err := utils.GenerateJWT(user.ID, secretKey)
	if err != nil {
		return "", fmt.Errorf("не удалось создать токен: %w", err)
	}
	return token, nil
}
