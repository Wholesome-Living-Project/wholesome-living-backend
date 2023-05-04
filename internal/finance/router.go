package finance

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	finance := app.Group("/investment")

	// add middlewares here

	// add routes here
	finance.Post("/", controller.create)
	finance.Get("/", controller.get)
	finance.Get("/:investmentID", controller.get)
}
