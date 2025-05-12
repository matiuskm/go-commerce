package models

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint
	ProductID uint
	Product   Product
	Quantity  int
}