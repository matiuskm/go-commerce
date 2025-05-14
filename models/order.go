package models

type Order struct {
	ID        	uint    		`gorm:"primaryKey" json:"id"`
	OrderNum  	string  		`gorm:"uniqueIndex" json:"orderNum"`
	UserID    	uint    		`json:"userId"`
	User	  	User			`json:"user"`
	Status	  	string 			`gorm:"default:'pending'"`    
	Items 	  	[]OrderItem
	Total     	int     		`gorm:"not null" json:"total"`
}

type OrderResponse struct {
	ID        	uint    			`json:"id"`
	OrderNum  	string  			`json:"orderNum"`
	Status	  	string 				`json:"status"`    
	Items 	  	[]OrderItemResponse	`json:"items"`
	Total     	int     			`json:"total"`
}