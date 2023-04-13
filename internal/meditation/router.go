package meditation

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	meditation := app.Group("/meditation")

	// add middlewares here

	// add routes here
	meditation.Post("/", controller.create)
	meditation.Get("/getAll/:userID", controller.getAll)
	meditation.Get("/:meditationID", controller.get)
}
