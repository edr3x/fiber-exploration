package model

type CreateUserInput struct {
	Name     string `json:"name" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}
