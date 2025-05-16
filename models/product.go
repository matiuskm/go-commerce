package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `gorm:"unique;not null" json:"name"`
	Description string  `gorm:"not null" json:"description"`
	Price       int 	`gorm:"not null" json:"price"`
	Stock       int  	`gorm:"not null" json:"stock"`
	ImageURL    string  `json:"image_url"`
}

type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       int 	`json:"price"`
	Stock       int  	`json:"stock"`
	ImageURL    string  `json:"image_url"`
}