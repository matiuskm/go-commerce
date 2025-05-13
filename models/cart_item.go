package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	CartID uint `json:"cartId"`
	ProductID uint `json:"productId"`
	Product Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Qty int `json:"quantity"`
}