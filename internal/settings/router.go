package settings

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	setteings := app.Group("/settings")

	// add middlewares here

	// add routes here
	setteings.Post("/", controller.createOnboarding)
	setteings.Get("/", controller.get)
	// create a route for each plugin
	setteings.Post("/finance", controller.createFinanceSettings)
	setteings.Post("/meditation", controller.createMeditationSettings)
	// Put for each plugin
	setteings.Put("/finance", controller.UpdateFinanceSettings)
	setteings.Put("/meditation", controller.updateMeditationSettings)

	setteings.Delete("/", controller.delete)
}
