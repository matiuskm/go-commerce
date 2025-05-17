package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string  `gorm:"not null" json:"name"`
	Username string  `gorm:"unique;not null" json:"username"`
	Password string  `gorm:"not null" json:"password"`
	Email    string  `gorm:"unique;not null" json:"email"`
	Role     string  `gorm:"not null" json:"role"` // "admin" or "user"
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}