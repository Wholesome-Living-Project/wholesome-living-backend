package finance

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	finance := app.Group("/investment")

	// add middlewares here

	// add routes here
	finance.Post("/", controller.create)
	finance.Get("/:userId", controller.get)
}
