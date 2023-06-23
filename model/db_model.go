package model

import "gorm.io/gorm"

type Role string

const (
	Admin    Role = "admin"
	Customer Role = "customer"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"username" gorm:"unique"`
	Role     Role   `json:"role" gorm:"default:customer"`
	Password string `json:"password"`
}
