package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edr3x/fiber-explore/config"
	"github.com/edr3x/fiber-explore/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matthewhartstonge/argon2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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
		"id":   user.ID,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
		"role": user.Role,
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

type ProfileDetailsResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func GetProfileDetailsService(userId string) (ProfileDetailsResponse, error) {
	var user model.User

	ctx := context.Background()

	redisVal, err := config.Redis.Get(ctx, userId).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("provided key does not exist")
		}
	}

	if redisVal != "" {
		log.Println("got data from redis")
		if err := json.Unmarshal([]byte(redisVal), &user); err != nil {
			log.Println("error unmarshalling redis data")
		}
	} else {
		log.Println("redis hit but didn't get data")
		if res := config.DB.First(&user, "id = ?", userId); res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				return ProfileDetailsResponse{}, fiber.NewError(404, "user not found")
			}
			return ProfileDetailsResponse{}, res.Error
		}
		value, error := json.Marshal(user)
		if error != nil {
			log.Println("error marshalling redis data")
		}

		if err := config.Redis.Set(ctx, userId, value, 40*time.Minute).Err(); err != nil {
			log.Println("redis set error")
		}
	}

	userResponse := ProfileDetailsResponse{
		Id:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Age:   user.Age,
	}

	return userResponse, nil
}

type UserServiceResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Role  string `json:"role"`
}

func GetAllUsersService() ([]UserServiceResponse, error) {
	var users []model.User

	if dbres := config.DB.Find(&users); dbres.Error != nil {
		return nil, dbres.Error
	}

	var usersRes []UserServiceResponse

	for _, user := range users {
		usersRes = append(usersRes, UserServiceResponse{
			Email: user.Email,
			Name:  user.Name,
			Age:   user.Age,
			Role:  string(user.Role),
		})
	}

	return usersRes, nil
}
