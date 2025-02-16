package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"merch-shop/internal/app/models"
)

type InventoryRepositoryInterface interface {
	CreateOrUpdate(tx *gorm.DB, userID uint, merchID uint) error
	GetInventoryByUser(tx *gorm.DB, userID uint) ([]models.Inventory, error)
}

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepositoryInterface {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) CreateOrUpdate(tx *gorm.DB, userID uint, merchID uint) error {
	if tx == nil {
		tx = r.db
	}
	var inventory models.Inventory
	err := tx.Where("user_id = ? AND merch_id = ?", userID, merchID).First(&inventory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			inventory = models.Inventory{
				UserID:   userID,
				MerchID:  merchID,
				Quantity: 1,
			}
			return tx.Create(&inventory).Error
		}
		return err
	}

	inventory.Quantity++
	return tx.Save(&inventory).Error
}

func (r *InventoryRepository) GetInventoryByUser(tx *gorm.DB, userID uint) ([]models.Inventory, error) {
	if tx == nil {
		tx = r.db
	}

	var inventory []models.Inventory

	err := tx.Where("user_id = ?", userID).Find(&inventory).Error
	if err != nil {
		return nil, fmt.Errorf("не удалось получить мерч пользователя: %w", err)
	}

	if len(inventory) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return inventory, nil
}
