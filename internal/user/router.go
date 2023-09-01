package user

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	user := app.Group("/users")

	// add middlewares here

	// add routes here
	user.Post("/", controller.create)
	user.Put("/", controller.update)
	user.Get("/", controller.getAll)
	user.Get("/:id", controller.get)
	user.Delete("/:id", controller.delete)
}
