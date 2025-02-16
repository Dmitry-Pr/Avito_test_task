// Package services Description: Описывает сервис для работы с мерчем.
package services

import (
	"errors"
	"fmt"
	"merch-shop/internal/app/models"
	"merch-shop/internal/app/repositories"

	"gorm.io/gorm"
)

// MerchServiceInterface описывает сервис для работы с мерчем.
type MerchServiceInterface interface {
	BuyMerch(userID uint, merchName string) error
}

// MerchService сервис для работы с мерчем.
type MerchService struct {
	repo            repositories.MerchRepositoryInterface
	userRepo        repositories.UserRepositoryInterface
	transactionRepo repositories.TransactionRepositoryInterface
	inventoryRepo   repositories.InventoryRepositoryInterface
}

// NewMerchService создает новый сервис для работы с мерчем.
func NewMerchService(
	repo repositories.MerchRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	transactionRepo repositories.TransactionRepositoryInterface,
	inventoryRepo repositories.InventoryRepositoryInterface,
) MerchServiceInterface {
	return &MerchService{repo: repo, userRepo: userRepo, transactionRepo: transactionRepo, inventoryRepo: inventoryRepo}
}

// BuyMerch покупает мерч.
func (s *MerchService) BuyMerch(userID uint, merchName string) error {
	return s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		merch, err := s.repo.GetMerchByName(tx, merchName)
		if err != nil {
			return fmt.Errorf("мерч не найден: %w", err)
		}

		user, err := s.userRepo.FindByID(tx, userID)
		if err != nil {
			return fmt.Errorf("пользователь не найден: %w", err)
		}

		if user.Coins < merch.Price {
			return errors.New("недостаточно монет")
		}

		if err := s.transactionRepo.Create(tx, &models.Transaction{
			UserID: userID,
			Type:   repositories.BuyType,
			Amount: -merch.Price,
		}); err != nil {
			return fmt.Errorf("не удалось создать транзакцию: %w", err)
		}

		user.Coins -= merch.Price
		if err := s.userRepo.Save(tx, user); err != nil {
			return fmt.Errorf("не удалось обновить количество монет пользователя: %w", err)
		}

		merchID := merch.ID
		if err := s.inventoryRepo.CreateOrUpdate(tx, userID, merchID); err != nil {
			return fmt.Errorf("не удалось добавить мерч пользователю: %w", err)
		}

		return nil
	})
}
