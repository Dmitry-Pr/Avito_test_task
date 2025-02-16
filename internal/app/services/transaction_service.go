package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"merch-shop/internal/app/models"
	"merch-shop/internal/app/repositories"
)

type TransactionServiceInterface interface {
	GetUserTransactionInfo(userID uint) (interface{}, error)
	SendCoins(fromUserID uint, toUsername string, amount int) error
}

type TransactionService struct {
	userRepo        repositories.UserRepositoryInterface
	transactionRepo repositories.TransactionRepositoryInterface
	merchRepo       repositories.MerchRepositoryInterface
	inventoryRepo   repositories.InventoryRepositoryInterface
}

func NewTransactionService(
	userRepo repositories.UserRepositoryInterface,
	transactionRepo repositories.TransactionRepositoryInterface,
	merchRepo repositories.MerchRepositoryInterface,
	inventoryRepo repositories.InventoryRepositoryInterface,
) TransactionServiceInterface {
	return &TransactionService{userRepo: userRepo, transactionRepo: transactionRepo, merchRepo: merchRepo, inventoryRepo: inventoryRepo}
}

func (s *TransactionService) GetUserTransactionInfo(userID uint) (interface{}, error) {
	transactions, err := s.transactionRepo.GetTransactionsByUser(nil, userID)
	if err != nil {
		return nil, fmt.Errorf("транзакции не найдены: %w", err)
	}

	user, err := s.userRepo.FindByID(nil, userID)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}

	inventory := make([]map[string]interface{}, 0)

	userInventory, err := s.inventoryRepo.GetInventoryByUser(nil, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("не удалось получить мерч пользователя: %w", err)
	}

	merchMap, err := s.merchRepo.GetAll(nil)
	if err != nil {
		return nil, fmt.Errorf("мерч не найден: %w", err)
	}

	for _, item := range userInventory {
		inventory = append(inventory, map[string]interface{}{
			"type":     merchMap[item.MerchID],
			"quantity": item.Quantity,
		})
	}

	received := make([]map[string]interface{}, 0)
	sent := make([]map[string]interface{}, 0)

	for _, t := range transactions {
		if t.Type == repositories.TransferType && *t.ToUserID == userID {
			fromUser, err := s.userRepo.FindByID(nil, *t.FromUserID)
			if err != nil {
				return nil, fmt.Errorf("отправитель не найден: %w", err)
			}
			received = append(received, map[string]interface{}{
				"fromUser": fromUser.Username,
				"amount":   t.Amount,
			})
		} else if t.Type == repositories.TransferType && *t.FromUserID == userID {
			toUser, err := s.userRepo.FindByID(nil, *t.ToUserID)
			if err != nil {
				return nil, fmt.Errorf("получатель не найден: %w", err)
			}
			sent = append(sent, map[string]interface{}{
				"toUser": toUser.Username,
				"amount": t.Amount,
			})
		}
	}

	return map[string]interface{}{
		"coins":     user.Coins,
		"inventory": inventory,
		"coinHistory": map[string]interface{}{
			"received": received,
			"sent":     sent,
		},
	}, nil
}

func (s *TransactionService) SendCoins(fromUserID uint, toUsername string, amount int) error {
	return s.merchRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		fromUser, err := s.userRepo.FindByID(tx, fromUserID)
		if err != nil {
			return fmt.Errorf("пользователь отправитель не найден: %w", err)
		}

		if fromUser.Username == toUsername {
			return fmt.Errorf("нельзя отправить монеты себе")
		}

		toUser, err := s.userRepo.FindByUsername(tx, toUsername)
		if err != nil {
			return fmt.Errorf("пользователь получатель не найден: %w", err)
		}

		if fromUser.Coins < amount {
			return errors.New("недостаточно монет")
		}

		if err := s.transactionRepo.Create(tx, &models.Transaction{
			FromUserID: &fromUser.ID,
			UserID:     fromUser.ID,
			ToUserID:   &toUser.ID,
			Type:       "transfer",
			Amount:     amount,
		}); err != nil {
			return fmt.Errorf("не удалось передать монеты: %w", err)
		}

		fromUser.Coins -= amount
		if err := s.userRepo.Save(tx, fromUser); err != nil {
			return fmt.Errorf("не удалось обновить баланс отправителя: %w", err)
		}

		toUser.Coins += amount
		if err := s.userRepo.Save(tx, toUser); err != nil {
			return fmt.Errorf("не удалось обновить баланс получателя: %w", err)
		}
		return nil
	})
}
