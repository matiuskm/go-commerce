package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/helpers"
	"github.com/matiuskm/go-commerce/models"
)

func CheckoutHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if (!exists) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}
	userID := userIDAny.(uint)

	var cart models.Cart
	if err := db.DB.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch cart"})
		return
	}

	var total int
	var orderItems []models.OrderItem

	tx := db.DB.Begin()

	for _, item := range cart.Items {
		product := item.Product

		if product.Stock < item.Qty {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Insufficient stock for product %s", product.Name)})
			return
		}

		// reduce stock
		product.Stock -= item.Qty
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}

		// add to orderItems
		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Quantity: item.Qty,
		})

		total += product.Price * item.Qty
	}

	// save order
	order := models.Order{
		OrderNum: helpers.GenerateOrderNumber(),
		UserID: userID,
		Total: total,
		Items: orderItems,
	}

	if err := tx.Create(&order).Error; err!= nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// clear cart
	if err := tx.Where("CartID =?", cart.ID).Delete(&models.CartItem{}).Error; err!= nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart items"})
		return
	}

	if err := tx.Delete(&cart).Error; err!= nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Checkout success", "order": order.OrderNum})
}