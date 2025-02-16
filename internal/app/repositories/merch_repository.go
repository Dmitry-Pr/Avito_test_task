package repositories

import (
	"merch-shop/internal/app/models"

	"gorm.io/gorm"
)

type Merch struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

type MerchRepositoryInterface interface {
	GetAll(tx *gorm.DB) (map[uint]string, error)
	GetMerchByName(tx *gorm.DB, name string) (*models.Merch, error)
	GetDB() *gorm.DB
}

type MerchRepository struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) MerchRepositoryInterface {
	return &MerchRepository{db: db}
}

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

func (r *MerchRepository) GetDB() *gorm.DB {
	return r.db
}
