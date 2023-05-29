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
	setteings.Post("/elevator", controller.createElevatorSettings)

	// Put for each plugin
	setteings.Put("/finance", controller.updateFinanceSettings)
	setteings.Put("/meditation", controller.updateMeditationSettings)
	setteings.Put("/elevator", controller.updateElevatorSettings)

	setteings.Delete("/", controller.delete)
}
