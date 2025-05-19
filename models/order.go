package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderNum  		string  			`gorm:"uniqueIndex" json:"orderNum"`
	UserID    		uint    			`json:"userId"`
	User	  		User				`json:"user"`
	AddressID 		*uint    			`json:"addressId"`
	Address	  		Address				`json:"address"`
	Status	  		string 				`gorm:"default:'pending'"`    
	Items 	  		[]OrderItem
	Total     		int     			`gorm:"not null" json:"total"`
	XenditInvoice	string				`json:"xenditInvoice"`
	XenditUrl		string				`json:"xenditUrl"`
}

type OrderResponse struct {
	ID        		uint    			`json:"id"`
	OrderNum  		string  			`json:"orderNum"`
	Status	  		string 				`json:"status"`    
	Items 	  		[]OrderItemResponse	`json:"items"`
	Total     		int     			`json:"total"`
	Address	  		AddressResponse		`json:"address"`
	CreatedAt 		string				`json:"createdAt"`
}

type AdminOrderResponse struct {
	ID        		uint    			`json:"id"`
	OrderNum  		string  			`json:"orderNum"`
	Status	  		string 				`json:"status"`
	User 			UserResponse		`json:"user"`   
	Items 	  		[]OrderItemResponse	`json:"items"`
	Total     		int     			`json:"total"`
	Address	  		AddressResponse		`json:"address"`
	CreatedAt 		string				`json:"createdAt"`
	XenditInvoice	string				`json:"xenditInvoice"`
	XenditUrl		string				`json:"xenditUrl"`
}	