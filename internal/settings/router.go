package settings

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	setteings := app.Group("/settings")

	// add middlewares here

	// add routes here
	setteings.Post("/", controller.createOnboarding)
	setteings.Get("/", controller.get)
}
