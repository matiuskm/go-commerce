package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/models"
)

type CartPayload struct {
	Items []models.CartItem `json:"items"`
}

type SimplifiedCartItem struct {
	ProductID uint            `json:"product_id"`
	Product   models.Product  `json:"product"`
	Quantity  int             `json:"quantity"`
}

func SaveCartHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	userID := userIDAny.(uint)

	var payload CartPayload
	if err := c.ShouldBindJSON(&payload); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cart models.Cart
	db.DB.Where("user_id = ?", userID).FirstOrInit(&cart)
	cart.UserID = &userID
	db.DB.Save(&cart)

	// clear old items
	db.DB.Where("cart_id =?", cart.ID).Delete(&models.CartItem{})

	// Deduplicate items based on ProductID
	mergedItems := map[uint]int{}
	for _, item := range payload.Items {
		mergedItems[item.ProductID] += item.Qty
	}

	for productID, quantity := range mergedItems {
		item := models.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			Qty:  quantity,
		}
		db.DB.Create(&item)
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart saved"})
}

func GetSavedCartHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if!exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	userID := userIDAny.(uint)
	
	var cart models.Cart
	if err := db.DB.Preload("Items.Product").Where("user_id =?", userID).First(&cart).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"items": []SimplifiedCartItem{}})
		return
	}

	simplified := []SimplifiedCartItem{}
	for _, item := range cart.Items {
		simplified = append(simplified, SimplifiedCartItem{
			ProductID: item.ProductID,
			Product:   item.Product,
			Quantity:  item.Qty,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": simplified})
}