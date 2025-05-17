package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/models"
)

func GetOrderHistoryHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if (!exists) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised"})
		return
	}
	userID := userIDAny.(uint)

	var orders []models.Order
	if err := db.DB.Where("user_id = ?", userID).Preload("Items.Product").Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch order history"})
		return
	}

	var response []models.OrderResponse
	for _, o := range orders {
		order := models.OrderResponse{
			ID: o.ID,
			OrderNum: o.OrderNum,
			Status: o.Status,
			Total: o.Total,
			CreatedAt: o.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		for _, item := range o.Items {
			simpleProduct := models.ProductResponse{
				ID: item.Product.ID,
				Name: item.Product.Name,
				Price: item.Product.Price,
				Description: item.Product.Description,
				Stock: item.Product.Stock,
			}
	
			order.Items = append(order.Items, models.OrderItemResponse{
				Product: simpleProduct,
				Quantity: item.Quantity,
			})
		}

		response = append(response, order)
	}

	c.JSON(http.StatusOK, gin.H{"orders": response})
}

func GetOrderDetailHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if (!exists) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised"})
		return
	}
	userID := userIDAny.(uint)
	
	orderID := c.Param("id")
	
	var order models.Order
	if err := db.DB.
		Preload("Address").
		Preload("Items.Product").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	address := models.AddressResponse{
		Label: order.Address.Label,
		Phone: order.Address.Phone,
		Street: order.Address.Street,
		RecipientName: order.Address.RecipientName,
	}

	response := models.OrderResponse{
		ID: order.ID,
		OrderNum: order.OrderNum,
		Status: order.Status,
		Total: order.Total,
		Address: address,
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	for _, item := range order.Items {
		simpleProduct := models.ProductResponse{
			ID: item.Product.ID,
			Name: item.Product.Name,
			Price: item.Product.Price,
			Description: item.Product.Description,
			Stock: item.Product.Stock,
			ImageURL: item.Product.ImageURL,
		}

		response.Items = append(response.Items, models.OrderItemResponse{
			Product: simpleProduct,
			Quantity: item.Quantity,
		})
	}


	c.JSON(http.StatusOK, gin.H{"order": response})
}