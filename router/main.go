package router

import (
	"github.com/Wholesome-Living-Project/wholesome-living-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", handlers.HandleHealthCheck)

	// setup the todos group
	todos := app.Group("/todos")
	todos.Get("/", handlers.HandleAllTodos)
	todos.Post("/", handlers.HandleCreateTodo)
	todos.Put("/:id", handlers.HandleUpdateTodo)
	todos.Get("/:id", handlers.HandleGetOneTodo)
	todos.Delete("/:id", handlers.HandleDeleteTodo)

	// user management
	user := app.Group("/user")
	user.Post("/", handlers.HandleCreateUser)
	user.Get("/:id", handlers.HandleGetUser)

	// meditation management
	meditation := app.Group("/meditation")
	meditation.Post("/:id", handlers.HandleCreateMeditation)
	meditation.Get("/:id", handlers.HandleGetMeditation)
}
