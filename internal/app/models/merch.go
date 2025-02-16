// Package models Description: This file contains the merch model.
package models

import "gorm.io/gorm"

// Merch model
type Merch struct {
	gorm.Model
	Name  string
	Price int
}
