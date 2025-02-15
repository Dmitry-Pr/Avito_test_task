package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Coins    int `gorm:"default:1000;check:coins >= 0"`
}
