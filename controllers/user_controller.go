package controller

import (
	"github.com/edr3x/fiber-explore/middlewares"
	"github.com/edr3x/fiber-explore/model"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	var body model.CreateUserInput

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.FailureResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	if errors := middlewares.ValidateInput(&body); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.FailureResponse{
			Success: false,
			Message: errors,
		})
	}

	return c.Status(201).JSON(model.Response{
		Success: true,
		Payload: body,
	})
}

func GetUser(c *fiber.Ctx) error {
	return c.Status(200).JSON(model.Response{
		Success: true,
		Payload: "User route",
	})
}
