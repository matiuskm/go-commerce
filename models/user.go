package models

type User struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Name     string  `gorm:"not null" json:"name"`
	Username string  `gorm:"unique;not null" json:"username"`
	Password string  `gorm:"not null" json:"-"`
	Email    string  `gorm:"unique;not null" json:"email"`
	Role     string  `gorm:"not null" json:"role"` // "admin" or "user"
}