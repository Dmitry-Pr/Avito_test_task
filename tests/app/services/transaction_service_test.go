package services_test

import (
	"errors"
	"testing"

	"gorm.io/gorm"

	"merch-shop/internal/app/models"
	"merch-shop/internal/app/services"
	mockrepositories "merch-shop/mocks/repositories" // Импортируем сгенерированные моки

	"github.com/golang/mock/gomock" // Импортируем gomock
	"github.com/stretchr/testify/assert"
)

func ptrUint(i uint) *uint {
	return &i
}

func TestGetUserTransactionInfo(t *testing.T) {
	testCases := []struct {
		name              string
		user              *models.User
		transactions      []models.Transaction
		inventory         []models.Inventory
		merchMap          map[uint]string
		expectedErr       string
		expectedCoins     int
		expectedInventory []map[string]interface{}
		expectedReceived  []map[string]interface{}
		expectedSent      []map[string]interface{}
	}{
		{
			name: "Success",
			user: &models.User{Model: gorm.Model{ID: 1}, Coins: 100, Username: "testuser"},
			transactions: []models.Transaction{
				{FromUserID: ptrUint(2), UserID: 1, ToUserID: ptrUint(1), Type: "transfer", Amount: 50},
			},
			inventory:         []models.Inventory{{UserID: 1, MerchID: 1, Quantity: 2}},
			merchMap:          map[uint]string{1: "t-shirt"},
			expectedCoins:     100,
			expectedInventory: []map[string]interface{}{{"type": "t-shirt", "quantity": uint(2)}},
			expectedReceived:  []map[string]interface{}{{"fromUser": "testuser2", "amount": 50}}, // "testuser2" from mock
			expectedSent:      []map[string]interface{}{},
		},
		{
			name:        "User Not Found",
			user:        nil,
			expectedErr: "пользователь не найден",
		},
		{
			name:              "No Transactions", // Случай, когда транзакций нет вообще
			user:              &models.User{Model: gorm.Model{ID: 1}, Coins: 100, Username: "testuser"},
			transactions:      []models.Transaction{}, // Пустой слайс
			inventory:         []models.Inventory{{UserID: 1, MerchID: 1, Quantity: 2}},
			merchMap:          map[uint]string{1: "t-shirt"},
			expectedCoins:     100,
			expectedInventory: []map[string]interface{}{{"type": "t-shirt", "quantity": uint(2)}},
			expectedReceived:  []map[string]interface{}{}, // Нет полученных транзакций
			expectedSent:      []map[string]interface{}{}, // Нет отправленных транзакций
		},
		{
			name: "Multiple Transactions", // Несколько транзакций
			user: &models.User{Model: gorm.Model{ID: 1}, Coins: 100, Username: "testuser"},
			transactions: []models.Transaction{
				{FromUserID: ptrUint(2), UserID: 1, ToUserID: ptrUint(1), Type: "transfer", Amount: 50},
				{FromUserID: ptrUint(1), UserID: 2, ToUserID: ptrUint(2), Type: "transfer", Amount: 25}, // Отправил
			},
			inventory:         []models.Inventory{{UserID: 1, MerchID: 1, Quantity: 2}},
			merchMap:          map[uint]string{1: "t-shirt"},
			expectedCoins:     100,
			expectedInventory: []map[string]interface{}{{"type": "t-shirt", "quantity": uint(2)}},
			expectedReceived:  []map[string]interface{}{{"fromUser": "testuser2", "amount": 50}},
			expectedSent:      []map[string]interface{}{{"toUser": "testuser2", "amount": 25}}, // "testuser2" from mock
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mockrepositories.NewMockUserRepositoryInterface(ctrl)
			transactionRepo := mockrepositories.NewMockTransactionRepositoryInterface(ctrl)
			merchRepo := mockrepositories.NewMockMerchRepositoryInterface(ctrl)
			inventoryRepo := mockrepositories.NewMockInventoryRepositoryInterface(ctrl)
			service := services.NewTransactionService(userRepo, transactionRepo, merchRepo, inventoryRepo)

			if tc.user != nil {
				userRepo.EXPECT().FindByID(nil, tc.user.ID).Return(tc.user, nil).AnyTimes()

				user2 := &models.User{Model: gorm.Model{ID: 2}, Coins: 50, Username: "testuser2"} // Создаем пользователя 2 для мока
				userRepo.EXPECT().FindByID(nil, user2.ID).Return(user2, nil).AnyTimes()           // Ожидание для пользователя 2
			} else {
				userRepo.EXPECT().FindByID(nil, uint(1)).Return(nil, errors.New(tc.expectedErr))
			}

			if tc.user != nil {
				transactionRepo.EXPECT().GetTransactionsByUser(nil, tc.user.ID).Return(tc.transactions, nil)
			}

			if tc.user != nil {
				inventoryRepo.EXPECT().GetInventoryByUser(nil, tc.user.ID).Return(tc.inventory, nil)
			}

			if tc.user != nil {
				merchRepo.EXPECT().GetAll(nil).Return(tc.merchMap, nil)
			}

			if tc.user == nil {
				tc.user = &models.User{Model: gorm.Model{ID: 1}, Coins: 100, Username: "testuser"}
			}

			result, err := service.GetUserTransactionInfo(tc.user.ID)

			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				infoMap := result.(map[string]interface{})
				assert.Equal(t, tc.expectedCoins, infoMap["coins"])

				inv, ok := infoMap["inventory"].([]map[string]interface{})
				assert.True(t, ok, "inventory type assertion failed")
				assert.Equal(t, tc.expectedInventory, inv)

				coinHistory, ok := infoMap["coinHistory"].(map[string]interface{})
				assert.True(t, ok, "coinHistory type assertion failed")

				received, ok := coinHistory["received"].([]map[string]interface{})
				assert.True(t, ok, "received type assertion failed")
				assert.Equal(t, tc.expectedReceived, received)

				sent, ok := coinHistory["sent"].([]map[string]interface{})
				assert.True(t, ok, "sent type assertion failed")
				assert.Equal(t, tc.expectedSent, sent)
			}
		})
	}
}
