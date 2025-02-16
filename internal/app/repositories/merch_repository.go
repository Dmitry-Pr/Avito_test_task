// Package repositories Description: этот файл содержит репозиторий для товаров.
package repositories

import (
	"merch-shop/internal/app/models"

	"gorm.io/gorm"
)

// Merch описывает товар.
type Merch struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

// MerchRepositoryInterface описывает репозиторий для товаров.
type MerchRepositoryInterface interface {
	GetAll(tx *gorm.DB) (map[uint]string, error)
	GetMerchByName(tx *gorm.DB, name string) (*models.Merch, error)
	GetDB() *gorm.DB
}

// MerchRepository репозиторий для товаров.
type MerchRepository struct {
	db *gorm.DB
}

// NewMerchRepository создает новый репозиторий для товаров.
func NewMerchRepository(db *gorm.DB) MerchRepositoryInterface {
	return &MerchRepository{db: db}
}

// GetAll получает все товары.
func (r *MerchRepository) GetAll(tx *gorm.DB) (map[uint]string, error) {
	if tx == nil {
		tx = r.db
	}
	merchMap := make(map[uint]string)
	var merchItems []Merch
	if err := tx.Model(&Merch{}).Find(&merchItems).Error; err != nil {
		return nil, err
	}
	for _, item := range merchItems {
		merchMap[item.ID] = item.Name
	}
	return merchMap, nil
}

// GetMerchByName получает товар по имени.
func (r *MerchRepository) GetMerchByName(tx *gorm.DB, name string) (*models.Merch, error) {
	if tx == nil {
		tx = r.db
	}
	var merch models.Merch
	if err := tx.Where("name = ?", name).First(&merch).Error; err != nil {
		return nil, err
	}
	return &merch, nil
}

// GetDB возвращает базу данных.
func (r *MerchRepository) GetDB() *gorm.DB {
	return r.db
}
