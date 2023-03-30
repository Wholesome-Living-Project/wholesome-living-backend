package meditation

import "github.com/gofiber/fiber/v2"

func MeditationRoutes(app *fiber.App, controller *MeditationController) {
	user := app.Group("/meditation")

	// add middlewares here

	// add routes here
	user.Post("/", controller.create)
	user.Get("/", controller.get)
}
