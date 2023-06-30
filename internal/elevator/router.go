package elevator

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	meditation := app.Group("/elevator")

	// add middlewares here

	// add routes here
	meditation.Post("/", controller.create)
	meditation.Get("/", controller.get)
}
