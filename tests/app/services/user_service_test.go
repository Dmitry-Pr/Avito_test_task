package services_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"merch-shop/internal/app/models"
	"merch-shop/internal/app/services"
	mockrepositories "merch-shop/mocks/repositories"
)

func TestAuthenticate(t *testing.T) {
	// Сохраняем оригинальное значение переменной окружения, чтобы восстановить его после теста
	originalSecretKey := os.Getenv("JWT_SECRET_KEY")
	defer func() {
		os.Setenv("JWT_SECRET_KEY", originalSecretKey)
	}()

	testCases := []struct {
		name        string
		username    string
		password    string
		mockUser    *models.User
		expectedErr string
	}{
		{
			name:     "Success - Existing User",
			username: "testuser",
			password: "password",
			mockUser: &models.User{Model: gorm.Model{ID: 1}, Username: "testuser", Password: "$2a$10$vFM2jouM/gav7wbKvm/tGuJg0EdgbMiXgcPUw350yJiCU5OFn5uYi"}, // Хэшированный пароль 			// Заглушка токена
		},
		{
			name:     "Success - New User",
			username: "newuser",
			password: "password",
			mockUser: nil,
		},
		{
			name:        "Invalid Password",
			username:    "testuser",
			password:    "wrongpassword",
			mockUser:    &models.User{Model: gorm.Model{ID: 1}, Username: "testuser", Password: "$2a$10$vFM2jouM/gav7wbKvm/tGuJg0EdgbMiXgcPUw350yJiCU5OFn5uYi"}, // Хэшированный пароль
			expectedErr: "неверные данные пользователя",
		},
		{
			name:        "Database Error - FindByUsername",
			username:    "testuser",
			password:    "password",
			mockUser:    nil,
			expectedErr: "ошибка базы данных", // Ошибка при поиске пользователя
		},
		{
			name:        "Database Error - Save",
			username:    "newuser",
			password:    "password",
			mockUser:    nil,
			expectedErr: "не удалось создать пользователя", // Ошибка при создании пользователя
		},
		{
			name:        "JWT Secret Key Not Found",
			username:    "testuser",
			password:    "password",
			mockUser:    &models.User{Model: gorm.Model{ID: 1}, Username: "testuser", Password: "$2a$10$vFM2jouM/gav7wbKvm/tGuJg0EdgbMiXgcPUw350yJiCU5OFn5uYi"}, // Хэшированный пароль
			expectedErr: "JWT_SECRET_KEY переменная среды не найдена",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockrepositories.NewMockUserRepositoryInterface(ctrl)
			service := services.NewUserService(repo)

			if tc.name == "JWT Secret Key Not Found" {
				os.Setenv("JWT_SECRET_KEY", "") // Устанавливаем пустое значение для теста
			} else {
				os.Setenv("JWT_SECRET_KEY", "test_secret_key") // Устанавливаем тестовый секретный ключ
			}

			if !(tc.mockUser == nil || tc.expectedErr == "ошибка базы данных" || tc.expectedErr == "не удалось создать пользователя") {
				repo.EXPECT().FindByUsername(nil, tc.username).Return(tc.mockUser, nil)
			} else {
				if tc.name == "Success - New User" || tc.expectedErr == "не удалось создать пользователя" {
					repo.EXPECT().FindByUsername(nil, tc.username).Return(nil, gorm.ErrRecordNotFound)
				} else {
					repo.EXPECT().FindByUsername(nil, tc.username).Return(nil, errors.New(tc.expectedErr))
				}
			}

			if tc.mockUser == nil && tc.expectedErr != "ошибка базы данных" {
				if tc.expectedErr == "не удалось создать пользователя" {
					repo.EXPECT().Save(nil, gomock.Any()).Return(fmt.Errorf(tc.expectedErr))
				} else {
					repo.EXPECT().Save(nil, gomock.Any()).Return(nil)
				}
			}

			token, err := service.Authenticate(tc.username, tc.password)

			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, "", token)
			}
		})
	}
}
