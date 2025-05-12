package models

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"unique;not null" json:"name"`
	Description string  `gorm:"not null" json:"description"`
	Price       int 	`gorm:"not null" json:"price"`
	Stock       int  	`gorm:"not null" json:"stock"`
}