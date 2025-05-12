package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/models"
)

func AdminGetUsersHandler(c *gin.Context) {}
func AdminGetOrdersHandler(c *gin.Context) {}

func AdminCreateProductHandler(c *gin.Context) {
	var req models.Product
	if err := c.ShouldBindJSON(&req); err!= nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&req).Error; err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "product": req})
}

func AdminUpdateProductHandler(c *gin.Context) {
	var req models.Product
	if err := c.ShouldBindJSON(&req); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	if err := db.DB.Where("id = ?", id).Updates(&req).Error; err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": req})
}

func AdminDeleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := db.DB.Where("id =?", id).Delete(&product).Error; err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete product"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Product deleted successfully"})
}
