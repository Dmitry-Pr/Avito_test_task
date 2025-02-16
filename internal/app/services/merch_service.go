package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"merch-shop/internal/app/models"
	"merch-shop/internal/app/repositories"
)

type MerchServiceInterface interface {
	GetAllMerch() ([]string, error)
	BuyMerch(userID uint, merchName string) error
}

type MerchService struct {
	repo            repositories.MerchRepositoryInterface
	userRepo        repositories.UserRepositoryInterface
	transactionRepo repositories.TransactionRepositoryInterface
}

func NewMerchService(
	repo repositories.MerchRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	transactionRepo repositories.TransactionRepositoryInterface,
) MerchServiceInterface {
	return &MerchService{repo: repo, userRepo: userRepo, transactionRepo: transactionRepo}
}

func (s *MerchService) GetAllMerch() ([]string, error) {
	return s.repo.GetAll(nil)
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

		return nil
	})
}
