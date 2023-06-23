package routes

import (
	"github.com/edr3x/fiber-explore/controllers"
	"github.com/edr3x/fiber-explore/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(router fiber.Router) {
	r := router.Group("/user")

	r.Get("/", middlewares.RequireAuth, controller.GetSelfDetails)

	r.Post("/", controller.CreateUser)

	r.Post("/login", controller.LoginController)
}
