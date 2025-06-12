package handlers

import (
	"io" // Ensure io is imported
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

	response := []models.UserResponse{}
	for _, user := range users {
		response = append(response, models.UserResponse{
			ID: user.ID,
			Name: user.Name,
			Username: user.Username,
			Email: user.Email,
			Role: user.Role,
		})
	}
	c.JSON(http.StatusOK, gin.H{"users": response})
}

func AdminGetUserHandler(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := db.DB.First(&user, id).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := models.UserResponse{
		ID: user.ID,
		Name: user.Name,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
	}

	c.JSON(http.StatusOK, gin.H{"user": response})
}

func AdminUpdateUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var payload struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Name = payload.Name
	user.Email = payload.Email
	db.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"user": user})
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
		user := models.UserResponse{
			ID: order.User.ID,
			Name: order.User.Name,
			Username: order.User.Username,
			Email: order.User.Email,
			Role: order.User.Role,
		}

		orderResponse := models.AdminOrderResponse{
			ID: order.ID,
			OrderNum: order.OrderNum,
			Status: order.Status,
			Total: order.Total,
			CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
			User: user,
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
	if err := db.DB.
		Preload("Address").
		Preload("User").
		Preload("Items.Product").
		First(&order, id).
		Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	address := models.AddressResponse{
		Label: order.Address.Label,
		Phone: order.Address.Phone,
		Street: order.Address.Street,
		RecipientName: order.Address.RecipientName,
	}

	user := models.UserResponse{
		ID: order.User.ID,
		Name: order.User.Name,
		Username: order.User.Username,
		Email: order.User.Email,
		Role: order.User.Role,
	}

	orderResponse := models.AdminOrderResponse{
		ID: order.ID,
		OrderNum: order.OrderNum,
		Status: order.Status,
		Total: order.Total,
		Address: address,
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		User: user,
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

	user := models.UserResponse{
		Name: order.User.Name,
		Username: order.User.Username,
		Email: order.User.Email,
		Role: order.User.Role,
	}

	orderResponse := models.AdminOrderResponse{
		ID: order.ID,
		OrderNum: order.OrderNum,
		Status: order.Status,
		Total: order.Total,
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		User: user,
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
	if err := db.DB.Unscoped().Order("ID asc").Find(&products).Error; err != nil {
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

type AdminProductUpdateRequest struct {
	Name        *string  `json:"name"`
	Price       *float64 `json:"price"`
	Description *string  `json:"description"`
	Stock       *int     `json:"stock"`
}

func AdminUpdateProductHandler(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req AdminProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Price != nil {
		product.Price = int(*req.Price)
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := db.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": product})
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

const maxUploadSize = 5 * 1024 * 1024 // 5 MB
var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

func AdminUploadImageHandler(c *gin.Context) {
	id := c.Param("id")

	// First, check if product exists to avoid unnecessary file processing
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image field is required"})
		return
	}

	// Validate file size
	if fileHeader.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds limit of 5MB"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer file.Close()

	// Read header for content type detection
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF { // Correctly check for io.EOF
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file for type detection"})
		return
	}

	contentType := http.DetectContentType(buffer)
	if !allowedMimeTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed types: JPEG, PNG, GIF."})
		return
	}

	// Reset file pointer to the beginning as Read moved it
	_, err = file.Seek(0, 0) // 0 means relative to start of file (io.SeekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset file read pointer"})
		return
	}

	imageURL, err := helpers.UploadToCloudinary(file, uuid.New().String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Cloudinary"})
		return
	}

	// Update product image URL
	product.ImageURL = imageURL
	if err := db.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product image URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "image_url": imageURL})
}