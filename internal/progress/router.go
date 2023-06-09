package progress

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	progress := app.Group("/progress")

	// add middlewares here

	// add routes here
	progress.Get("/", controller.get)
	// create a route for each plugin
}
