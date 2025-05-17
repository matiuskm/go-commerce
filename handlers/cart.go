package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/models"
)

type CartPayload struct {
	Items []models.CartItem `json:"items"`
}

func SaveCartHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	userID := userIDAny.(uint)

	var payload CartPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cart models.Cart
	db.DB.Where("user_id = ?", userID).FirstOrInit(&cart)
	cart.UserID = userID
	db.DB.Save(&cart)

	// DELETE ALL existing cart items (REPLACE MODE)
	db.DB.Where("cart_id = ?", cart.ID).Unscoped().Delete(&models.CartItem{})

	// Deduplicate + filter qty <= 0
	mergedItems := map[uint]int{}
	for _, item := range payload.Items {
		if item.Qty > 0 {
			mergedItems[item.ProductID] += item.Qty
		}
	}

	for productID, quantity := range mergedItems {
		item := models.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			Qty:       quantity,
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
		c.JSON(http.StatusNotFound, gin.H{"items": []models.CartItemResponse{}})
		return
	}

	response := models.CartResponse{
		ID: cart.ID,
	}

	for _, item := range cart.Items {
		simpleProduct := models.ProductResponse{
			ID: item.Product.ID,
			Name: item.Product.Name,
			Price: item.Product.Price,
			Description: item.Product.Description,
			Stock: item.Product.Stock,
			ImageURL: item.Product.ImageURL,
		}
		response.Items = append(response.Items, models.CartItemResponse{
			Product:   simpleProduct,
			Qty:  item.Qty,
		})
	}

	c.JSON(http.StatusOK, gin.H{"cart": response})
}

func RemoveCartItemHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	userID := userIDAny.(uint)

	productIDStr := c.Param("productID")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var cart models.Cart
	if err := db.DB.Where("user_id =?", userID).First(&cart).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	if err := db.DB.Where("cart_id =? AND product_id =?", cart.ID, productID).Unscoped().Delete(&models.CartItem{}).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed"})
}