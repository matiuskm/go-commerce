package models

type Order struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	UserID    uint    
	User      User    
	Items 	  []OrderItem
	Total     int     `gorm:"not null" json:"total"`
}