package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID     uint
	User       User `gorm:"foreignKey:UserID"`
	Type       string
	Amount     int
	FromUserID *uint `gorm:"nullable:true"`
	FromUser   *User `gorm:"foreignKey:FromUserID;constraint:OnUpdate:CASCADE;nullable:true"`
	ToUserID   *uint `gorm:"nullable:true"`
	ToUser     *User `gorm:"foreignKey:ToUserID;constraint:OnUpdate:CASCADE;nullable:true"`
}
