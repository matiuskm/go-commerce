package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/models"
)

func GetAllProductsHandler(c *gin.Context) {
	var products []models.Product
	if err := db.DB.Order("ID asc").Find(&products).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func GetProductByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := db.DB.First(&product, id).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}
