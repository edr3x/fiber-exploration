package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}
