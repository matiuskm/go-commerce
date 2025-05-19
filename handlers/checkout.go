package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/helpers"
	"github.com/matiuskm/go-commerce/models"
)

type CheckoutPayload struct {
    AddressID *uint `json:"addressId"` // harus dikirim dari FE
}

func CheckoutHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if (!exists) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}
	userID := userIDAny.(uint)

	user := models.User{}
	if err := db.DB.First(&user, userID).Error; err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	var pay CheckoutPayload
    if err := c.ShouldBindJSON(&pay); err != nil || pay.AddressID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "addressId is required"})
        return
    }

	var ship models.Address
    if err := db.DB.
        Where("id = ? AND user_id = ?", *pay.AddressID, userID).
        First(&ship).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address"})
        return
    }

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
		AddressID: pay.AddressID,
	}

	if err := tx.Create(&order).Error; err!= nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	
	if err := helpers.CreateXenditInvoice(tx, &order, user.Email); err!= nil {
		log.Printf("❌ CreateXenditInvoice failed: %v\n", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}
	log.Println("✅ CreateXenditInvoice succeeded, invoice URL:", order.XenditUrl)
	
	// clear cart
	if err := tx.Unscoped().Delete(&cart).Error; err!= nil {
		tx.Rollback()
		log.Printf("❌ Failed to clear cart: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}
	
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Checkout success", "order": order.OrderNum, "paymentUrl": order.XenditUrl})
}