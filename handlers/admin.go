package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/helpers"
	"github.com/matiuskm/go-commerce/models"
)

type StatusUpdatePayload struct {
	Status string `json:"status"`
}

func AdminListUsersHandler(c *gin.Context) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func AdminGetUserHandler(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := db.DB.First(&user, id).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func AdminUpdateUserHandler(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if err := db.DB.Where("id =?", id).Updates(&req).Error; err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": req})
}

func AdminDeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.DB.Where("id =?", id).Delete(&user).Error; err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "User deleted successfully"})
}

func AdminUpdateUserRoleHandler(c *gin.Context) {
	var payload struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || (strings.ToLower(payload.Role) != "admin" && strings.ToLower(payload.Role) != "user") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	userID := c.Param("id")
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Role = payload.Role
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user role updated"})
}

func AdminListOrdersHandler(c *gin.Context) {
	var orders []models.Order
	if err := db.DB.Preload("User").Preload("Items.Product").Order("ID desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orders not found"})
		return
	}

	responses := []models.AdminOrderResponse{}
	for _, order := range orders {
		orderResponse := models.AdminOrderResponse{
			ID: order.ID,
			OrderNum: order.OrderNum,
			Status: order.Status,
			Total: order.Total,
			CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
			User: order.User,
		}

		for _, item := range order.Items {
			simpleProduct := models.ProductResponse{
				ID: item.Product.ID,
				Name: item.Product.Name,
				Price: item.Product.Price,
				Description: item.Product.Description,
				Stock: item.Product.Stock,
			}

			orderResponse.Items = append(orderResponse.Items, models.OrderItemResponse{
				Product: simpleProduct,
				Quantity: item.Quantity,
			})
		}

		responses = append(responses, orderResponse)
	}

	c.JSON(http.StatusOK, gin.H{"orders": responses})
}

func AdminGetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := db.DB.Preload("User").Preload("Items.Product").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	orderResponse := models.AdminOrderResponse{
		ID: order.ID,
		OrderNum: order.OrderNum,
		Status: order.Status,
		Total: order.Total,
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		User: order.User,
	}

	for _, item := range order.Items {
		simpleProduct := models.ProductResponse{
			ID: item.Product.ID,
			Name: item.Product.Name,
			Price: item.Product.Price,
			Description: item.Product.Description,
			Stock: item.Product.Stock,
		}

		orderResponse.Items = append(orderResponse.Items, models.OrderItemResponse{
			Product: simpleProduct,
			Quantity: item.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{"order": orderResponse})
}

func AdminUpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var payload StatusUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status not provided"})
		return
	}

	validStatuses := map[string]bool {
		"pending": true, "paid": true, "shipped": true, "delivered": true, "canceled": true,
	}
	if !validStatuses[strings.ToLower(payload.Status)] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	var order models.Order
	if err := db.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = strings.ToLower(payload.Status)
	if err := db.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update order status"})
		return
	}

	orderResponse := models.AdminOrderResponse{
		ID: order.ID,
		OrderNum: order.OrderNum,
		Status: order.Status,
		Total: order.Total,
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		User: order.User,
	}

	for _, item := range order.Items {
		simpleProduct := models.ProductResponse{
			ID: item.Product.ID,
			Name: item.Product.Name,
			Price: item.Product.Price,
			Description: item.Product.Description,
			Stock: item.Product.Stock,
		}

		orderResponse.Items = append(orderResponse.Items, models.OrderItemResponse{
			Product: simpleProduct,
			Quantity: item.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated", "order": orderResponse})
}

func AdminListProductHandler(c *gin.Context) {
	var products []models.Product
	if err := db.DB.Order("ID asc").Find(&products).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Products not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func AdminGetProductHandler(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := db.DB.First(&product, id).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	response := models.ProductResponse{
		ID: product.ID,
		Name: product.Name,
		Price: product.Price,
		Description: product.Description,
		Stock: product.Stock,
		ImageURL: product.ImageURL,
	}

	c.JSON(http.StatusOK, gin.H{"product": response})
}

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

func AdminUploadImageHandler(c *gin.Context) {
	id := c.Param("id")

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}

	file, _ := fileHeader.Open()
	defer file.Close()

	imageURL, err := helpers.UploadToCloudinary(file, uuid.New().String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	// update product image URL
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	product.ImageURL = imageURL
	db.DB.Save(&product)

	c.JSON(http.StatusOK, gin.H{"message": "image uploaded", "image_url": imageURL})
}