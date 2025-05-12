package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/helpers"
	"github.com/matiuskm/go-commerce/models"
)

type UserProfileResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null"`
	Password string `json:"password,omitempty"`
}

func GetUserProfileHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userIDAny).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := UserProfileResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
	}
	c.JSON(http.StatusOK, response)
}

func UpdateUserProfileHandler(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if!exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	var user models.User
	if err := db.DB.First(&user, userID).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var request  = UserUpdateRequest{}
	if err := c.ShouldBindJSON(&request); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Name = request.Name
	user.Username = request.Username
	user.Email = request.Email
	if (request.Password != "") {
		hashedPassword, err:= helpers.HashPassword(request.Password)
		if err!= nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hashedPassword
	}
	if err := db.DB.Save(&user).Error; err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile updated successfully"})
}

func GetUserOrdersHandler(c *gin.Context) {}
