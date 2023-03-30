package user

import "github.com/gofiber/fiber/v2"

func AddUserRoutes(app *fiber.App, controller *UserController) {
	user := app.Group("/users")

	// add middlewares here

	// add routes here
	user.Post("/", controller.create)
	user.Get("/", controller.getAll)
}
