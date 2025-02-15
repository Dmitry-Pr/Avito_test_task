package repositories

import (
	"gorm.io/gorm"
)

type Merch struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

type MerchRepositoryInterface interface {
	GetAll() ([]string, error)
}

type MerchRepository struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) MerchRepositoryInterface {
	return &MerchRepository{db: db}
}

func (r *MerchRepository) GetAll() ([]string, error) {
	var merchList []string
	var merchItems []Merch
	if err := r.db.Model(&Merch{}).Select("name").Find(&merchItems).Error; err != nil {
		return nil, err
	}
	for _, item := range merchItems {
		merchList = append(merchList, item.Name)
	}
	return merchList, nil
}
