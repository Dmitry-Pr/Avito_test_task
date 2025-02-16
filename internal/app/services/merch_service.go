package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"merch-shop/internal/app/models"
	"merch-shop/internal/app/repositories"
)

type MerchServiceInterface interface {
	BuyMerch(userID uint, merchName string) error
}

type MerchService struct {
	repo            repositories.MerchRepositoryInterface
	userRepo        repositories.UserRepositoryInterface
	transactionRepo repositories.TransactionRepositoryInterface
	inventoryRepo   repositories.InventoryRepositoryInterface
}

func NewMerchService(
	repo repositories.MerchRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	transactionRepo repositories.TransactionRepositoryInterface,
	inventoryRepo repositories.InventoryRepositoryInterface,
) MerchServiceInterface {
	return &MerchService{repo: repo, userRepo: userRepo, transactionRepo: transactionRepo, inventoryRepo: inventoryRepo}
}

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
			Type:   "buy",
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
