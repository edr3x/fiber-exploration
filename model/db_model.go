package model

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	Admin    Role = "admin"
	Customer Role = "customer"
)

type User struct {
	ID        string `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Email     string `json:"username" gorm:"unique"`
	Role      Role   `json:"role" gorm:"default:customer"`
	Password  string `json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
