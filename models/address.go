package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	UserID 			uint		`json:"userId"`
	Label 			string		`json:"label"`
	RecipientName 	string		`json:"recipientName"`
	Phone 			string		`json:"phone"`
	Street 			string		`json:"street"`
}

type AddressResponse struct {
	ID 				uint		`json:"id"`
	Label 			string		`json:"label"`
	RecipientName 	string		`json:"recipientName"`
	Phone 			string		`json:"phone"`
	Street 			string		`json:"street"`
}