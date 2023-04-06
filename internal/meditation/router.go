package meditation

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	meditation := app.Group("/meditation")

	// add middlewares here

	// add routes here
	meditation.Post("/", controller.create)
	meditation.Get("/", controller.get)
	meditation.Get("/", controller.get)
	meditation.Get("/:meditationID", controller.get)
}
