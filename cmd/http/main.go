package main

import (
	"cmd/http/main.go/config"
	"cmd/http/main.go/internal/meditation"
	"cmd/http/main.go/internal/storage"
	"cmd/http/main.go/internal/user"
	"cmd/http/main.go/pkg/shutdown"

	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title Wholesome Living Backend
// @version 0.1
// @description A backend for Wholesome Living written in Golang backend API using Fiber and MongoDB
// @contact.name Wholesome Living
// @license.name MIT
// @BasePath /
func main() {
	// setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// run the server
	cleanup, err := run(env)

	// run the cleanup after the server is terminated
	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// ensure the server is shutdown gracefully & app runs
	shutdown.Gracefully()

}

func run(env config.EnvVars) (func(), error) {
	app, cleanup, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	// start the server
	go func() {
		err := app.Listen("0.0.0.0:" + env.PORT)
		if err != nil {
			return
		}
	}()

	// return a function to close the server and database
	return func() {
		cleanup()
		err := app.Shutdown()
		if err != nil {
			return
		}
	}, nil
}

func buildServer(env config.EnvVars) (*fiber.App, func(), error) {
	// init the storage
	db, err := storage.BootstrapMongo(env.MONGODB_URI, env.MONGODB_NAME, 10*time.Second)
	if err != nil {
		return nil, nil, err
	}

	// create the fiber app
	app := fiber.New()

	// add middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// add health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	// add docs
	app.Get("/swagger/*", swagger.HandlerDefault)

	// create the user domain
	userStore := user.NewStorage(db)
	userController := user.NewController(userStore)
	user.Routes(app, userController)

	//create meditation domain
	meditationStore := meditation.NewStorage(db)
	meditationController := meditation.NewController(meditationStore)
	meditation.Routes(app, meditationController)

	return app, func() {
		err := storage.CloseMongo(db)
		if err != nil {
			return
		}
	}, nil
}
