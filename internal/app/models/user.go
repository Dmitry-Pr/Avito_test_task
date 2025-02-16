// Package models Description: User model.
package models

import "gorm.io/gorm"

// User model
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Coins    int `gorm:"default:1000;check:coins >= 0"`
}
