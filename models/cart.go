package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint `gorm:"not null" json:"userId"`
	User *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Items []CartItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CartResponse struct {
	ID 		uint 				`json:"id"`
	Items  	[]CartItemResponse 	`json:"items"`
}