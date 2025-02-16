// Package models Description: Этот файл содержит модель inventory.
package models

import "gorm.io/gorm"

// Inventory описывает инвентарь.
type Inventory struct {
	gorm.Model
	UserID   uint
	User     User `gorm:"foreignKey:UserID"`
	MerchID  uint
	Merch    Merch `gorm:"foreignKey:MerchID"`
	Quantity uint
}
