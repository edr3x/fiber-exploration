package services

import (
	"github.com/edr3x/fiber-explore/config"
	"github.com/edr3x/fiber-explore/model"
	"github.com/matthewhartstonge/argon2"
)

func CreateUserService(user_data model.CreateUserInput) (interface{}, error) {
	argon := argon2.DefaultConfig()

	hash, err := argon.HashEncoded([]byte(user_data.Password))
	if err != nil {
		return nil, err
	}

	user := model.User{
		Email:    user_data.Email,
		Password: string(hash),
		Age:      user_data.Age,
		Name:     user_data.Name,
	}

	if result := config.DB.Create(&user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
