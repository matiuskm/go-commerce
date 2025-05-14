package models

type OrderItem struct {
	ID        uint    	`gorm:"primaryKey" json:"id"`
	OrderID   uint		`json:"orderId"`
	ProductID uint		`json:"productId"`
	Product   Product	`json:"product"`
	Quantity  int		`json:"quantity"`
}

type OrderItemResponse struct {
	Product		ProductResponse		`json:"product"`
	Quantity	int			`json:"quantity"`
}