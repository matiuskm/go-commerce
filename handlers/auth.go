package handlers

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/helpers"
	"github.com/matiuskm/go-commerce/models"
)

var (
	// Allow 10 requests per minute, with a burst of 5.
	loginLimiter    = rate.NewLimiter(rate.Every(6*time.Second), 5) // Equivalent to 10/min
	registerLimiter = rate.NewLimiter(rate.Every(6*time.Second), 5) // Equivalent to 10/min
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func RegisterHandler(c *gin.Context) {
	if !registerLimiter.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many registration attempts. Please try again later."})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err:= helpers.HashPassword(req.Password)
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name: req.Name,
		Email: req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Role: "user",
	}

	if err = db.DB.Create(&user).Error; err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	token, err:= helpers.GenerateJWT(user)
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusCreated, LoginResponse{Token: token})
}

func LoginHandler(c *gin.Context) {
	if !loginLimiter.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many login attempts. Please try again later."})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err:= db.DB.Where("username = ?", req.Username).First(&user).Error; err!= nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !helpers.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err:= helpers.GenerateJWT(user)
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}