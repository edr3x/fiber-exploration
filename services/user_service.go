package services

import (
	"os"
	"time"

	"github.com/edr3x/fiber-explore/config"
	"github.com/edr3x/fiber-explore/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

type LoginResponse struct {
	Token string    `json:"token"`
	Data  fiber.Map `json:"data"`
}

func LoginService(login_creds model.LoginInput) (LoginResponse, error) {
	var user model.User

	result := config.DB.Where("email = ?", login_creds.Email).First(&user)

	if result.Error != nil {
		return LoginResponse{}, result.Error
	}

	ok, err := argon2.VerifyEncoded([]byte(login_creds.Password), []byte(user.Password))
	if !ok || err != nil {
		return LoginResponse{}, err
	}

	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := userToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Token: tokenString,
		Data: fiber.Map{
			"email": user.Email,
			"name":  user.Name,
			"age":   user.Age,
		},
	}, nil
}
