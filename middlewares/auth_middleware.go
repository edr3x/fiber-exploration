package middlewares

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/edr3x/fiber-explore/config"
	"github.com/edr3x/fiber-explore/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	if res := config.DB.First(&user, claims["id"]); res.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err_res)
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
