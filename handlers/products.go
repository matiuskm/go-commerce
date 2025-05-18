package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/models"
)

func GetAllProductsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	var prods []models.Product
	if err := db.DB.Order("id asc").Offset(offset).Limit(limit).Find(&prods).Error; err != nil {
		c.JSON(500, gin.H{"error": "db error"})
		return
	}

	var total int64
	db.DB.Model(&models.Product{}).Count(&total)
	hasMore := int64(offset+limit) < total

	c.JSON(200, gin.H{
		"products": prods,
		"hasMore":  hasMore,
	})
}

func GetProductByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}
