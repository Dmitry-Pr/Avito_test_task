package models

import "gorm.io/gorm"

type Merch struct {
	gorm.Model
	Name  string
	Price int64
}
