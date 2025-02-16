package models

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	UserID   uint
	User     User `gorm:"foreignKey:UserID"`
	MerchID  uint
	Merch    Merch `gorm:"foreignKey:MerchID"`
	Quantity uint
}
