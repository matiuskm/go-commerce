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
	if err := db.DB.Where("user_id = ?", userID).Preload("User").Preload("Items.Product").Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch order history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}