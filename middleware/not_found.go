package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/edr3x/fiber-explore/model"
)

func NotFound() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(404).JSON(model.FailureResponse{
			Success: false,
			Message: "Endpoint Not found",
		})
	}
}
