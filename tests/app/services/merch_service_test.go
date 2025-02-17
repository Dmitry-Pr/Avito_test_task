package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"merch-shop/internal/app/models"
	"merch-shop/internal/app/repositories"
	"merch-shop/internal/app/services"
)

func TestBuyMerch(t *testing.T) {
	testCases := []struct {
		name              string
		userID            uint
		merchName         string
		mockUser          *models.User
		mockMerch         *models.Merch
		expectedInventory *models.Inventory
		expectedErr       string
	}{
		{
			name:              "Success",
			userID:            1,
			merchName:         "t-shirt",
			mockUser:          &models.User{Model: gorm.Model{ID: 1}, Coins: 100},
			mockMerch:         &models.Merch{Model: gorm.Model{ID: 1}, Name: "t-shirt", Price: 80},
			expectedInventory: &models.Inventory{UserID: 1, MerchID: 1, Quantity: 1},
		},
		{
			name:        "Merch Not Found",
			userID:      1,
			merchName:   "t-shirt",
			mockUser:    &models.User{Model: gorm.Model{ID: 1}, Coins: 100},
			expectedErr: "мерч не найден",
		},
		{
			name:        "User Not Found",
			userID:      1,
			merchName:   "t-shirt",
			mockMerch:   &models.Merch{Model: gorm.Model{ID: 1}, Name: "t-shirt", Price: 80},
			expectedErr: "пользователь не найден",
		},
		{
			name:        "Insufficient Coins",
			userID:      1,
			merchName:   "t-shirt",
			mockUser:    &models.User{Model: gorm.Model{ID: 1}, Coins: 50},
			mockMerch:   &models.Merch{Model: gorm.Model{ID: 1}, Name: "t-shirt", Price: 80},
			expectedErr: "недостаточно монет",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				sqlDB, err := db.DB()
				if err != nil {
					t.Error(err)
				}
				sqlDB.Close()
			}()

			err = db.AutoMigrate(&models.User{}, &models.Merch{}, &models.Transaction{}, &models.Inventory{})
			if err != nil {
				t.Fatal(err)
			}
			db.Create(tc.mockMerch)
			db.Create(tc.mockUser)

			repo := repositories.NewMerchRepository(db)
			userRepo := repositories.NewUserRepository(db)
			transactionRepo := repositories.NewTransactionRepository(db)
			inventoryRepo := repositories.NewInventoryRepository(db)

			service := services.NewMerchService(repo, userRepo, transactionRepo, inventoryRepo)

			err = service.BuyMerch(tc.userID, tc.merchName)

			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
			} else {
				userModel, _ := userRepo.FindByID(nil, tc.userID)
				inventoryModels, _ := inventoryRepo.GetInventoryByUser(nil, tc.userID)
				inventoryModel := inventoryModels[0]
				assert.Equal(t, tc.mockUser.Coins-tc.mockMerch.Price, userModel.Coins)
				assert.Equal(t, tc.expectedInventory.MerchID, inventoryModel.MerchID)
				assert.Equal(t, tc.expectedInventory.Quantity, inventoryModel.Quantity)
				assert.Equal(t, tc.expectedInventory.UserID, inventoryModel.UserID)
				assert.NoError(t, err)
			}
		})
	}
}
