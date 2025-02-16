// Package repositories Description: Этот файл содержит репозиторий для инвентаря.
package repositories

import (
	"errors"
	"fmt"
	"merch-shop/internal/app/models"

	"gorm.io/gorm"
)

// InventoryRepositoryInterface описывает репозиторий для инвентаря.
type InventoryRepositoryInterface interface {
	CreateOrUpdate(tx *gorm.DB, userID uint, merchID uint) error
	GetInventoryByUser(tx *gorm.DB, userID uint) ([]models.Inventory, error)
}

// InventoryRepository репозиторий для инвентаря.
type InventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository создает новый репозиторий для инвентаря.
func NewInventoryRepository(db *gorm.DB) InventoryRepositoryInterface {
	return &InventoryRepository{db: db}
}

// CreateOrUpdate создает или обновляет инвентарь.
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

// GetInventoryByUser получает инвентарь по пользователю.
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
