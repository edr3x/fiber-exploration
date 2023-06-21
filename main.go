package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/edr3x/fiber-explore/middleware"
	"github.com/edr3x/fiber-explore/model"
	"github.com/edr3x/fiber-explore/routes"
)

func main() {
	app := fiber.New()

	// Recover from a panic thrown by any handler in the stack
	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,DELETE,PATCH",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(model.Response{
			Success: true,
			Payload: "Routecheck success",
		})
	})

	api := app.Group("/api/v1")

	route.UserRoute(api)

	app.Use(middleware.NotFound())
	log.Fatal(app.Listen(":8080"))
}
