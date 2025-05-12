package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID *uint `gorm:"not null" json:"userId"`
	User *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Items []CartItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}