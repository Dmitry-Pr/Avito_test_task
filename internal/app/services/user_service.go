package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"merch-shop/internal/app/models"
	"merch-shop/internal/app/repositories"
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
		if errors.Is(err, gorm.ErrRecordNotFound) { // Пользователь не найден, создаем нового
			newUser := &models.User{Username: username, Password: password} // + хеширование пароля!
			if err := s.repo.Save(newUser); err != nil {
				return "", fmt.Errorf("не удалось создать пользователя: %w", err)
			}
			user = newUser // Присваиваем user только что созданному пользователю
		} else {
			return "", err // Другая ошибка
		}
	}

	// Проверка пароля (нужно добавить хеширование)
	if user.Password != password {
		return "", errors.New("неверные данные пользователя")
	}

	return "mock-jwt-token", nil // Тут должен быть JWT
}
