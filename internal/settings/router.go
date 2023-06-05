package settings

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	settings := app.Group("/settings")

	// add middlewares here

	// add routes here
	settings.Post("/", controller.createOnboarding)
	settings.Get("/", controller.get)
	// create a route for each plugin
	settings.Post("/finance", controller.createFinanceSettings)
	settings.Post("/meditation", controller.createMeditationSettings)
	settings.Post("/elevator", controller.createElevatorSettings)

	// Put for each plugin
	settings.Put("/finance", controller.updateFinanceSettings)
	settings.Put("/meditation", controller.updateMeditationSettings)
	settings.Put("/elevator", controller.updateElevatorSettings)

	settings.Delete("/", controller.delete)
}
