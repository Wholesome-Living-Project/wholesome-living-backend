package meditation

import "github.com/gofiber/fiber/v2"

func MeditationRoutes(app *fiber.App, controller *MeditationController) {
	meditation := app.Group("/meditation")

	// add middlewares here

	// add routes here
	meditation.Post("/", controller.create)
	meditation.Get("/", controller.get)
}
