package meditation

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App, controller *Controller) {
	meditation := app.Group("/meditation")

	// add middlewares here

	// add routes here
	meditation.Post("/", controller.create)
<<<<<<< HEAD
	meditation.Get("/:meditationID", controller.get)
=======
	meditation.Get("/", controller.get)
	meditation.Get("/", controller.get)
>>>>>>> e842873335bfa1f980665bba05fec1bfa34b774b
}
