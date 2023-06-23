package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/edr3x/fiber-explore/config"
	"github.com/edr3x/fiber-explore/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RequireAuth(c *fiber.Ctx) error {
	headerVal := c.Get("Authorization")

	err_res := model.FailureResponse{
		Success: false,
		Message: "Unauthorized",
	}

	if headerVal == "" {
		log.Println("no auth header provided")
		return c.Status(fiber.StatusUnauthorized).JSON(err_res)
	}

	token, _ := jwt.Parse(strings.Split(headerVal, " ")[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(err_res)

	}
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		log.Println("token expired")
		return c.Status(fiber.StatusUnauthorized).JSON(err_res)
	}

	var user model.User
	userId := claims["id"].(string)

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
				return c.Status(fiber.StatusUnauthorized).JSON(model.FailureResponse{
					Success: false,
					Message: "User not found",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(err_res)
		}
		value, error := json.Marshal(user)
		if error != nil {
			log.Println("error marshalling redis data")
		}

		if err := config.Redis.Set(ctx, userId, value, 40*time.Minute).Err(); err != nil {
			log.Println("redis set error")
		}
	}

	c.Locals("user", user)

	return c.Next()
}

func RequireAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)
	if user.Role != model.Admin {
		return c.Status(fiber.StatusUnauthorized).JSON(model.FailureResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	return c.Next()
}
