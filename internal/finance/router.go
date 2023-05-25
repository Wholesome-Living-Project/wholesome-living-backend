package finance

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	finance := app.Group("/finance")

	// add middlewares here

	// add routes here
	finance.Post("/", controller.create)
	finance.Get("/", controller.get)
}
