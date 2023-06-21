package controller

import (
	"github.com/edr3x/fiber-explore/model"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	return c.Status(200).JSON(model.Response{
		Success: true,
		Payload: "User route",
	})
}
