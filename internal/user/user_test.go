package user

import (
	"testing"
	"time"

	"cmd/http/main.go/internal/storage"

	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// add Testify package

func TestGetUser(t *testing.T) {
	assert := assert.New(t)

	// Define a structure for specifying input and output data
	// of a single test case

	tests := []struct {
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
	}{
		// First test case
		{
			description:  "Health check",
			route:        "/health",
			expectedCode: 200,
		},
		{
			description:  "get HTTP status 200",
			route:        "/",
			expectedCode: 200,
		},
		// Second test case
		{
			description:  "get HTTP status 404, when route is not exists",
			route:        "/not-found",
			expectedCode: 404,
		},
	}

	// Define Fiber app.
	app := fiber.New()
	MONGODB_URI := "mongodb://localhost:27017"
	MONGODB_NAME := "testing"

	db, err := storage.BootstrapMongo(MONGODB_URI, MONGODB_NAME, 10*time.Second)

	if err != nil {
		t.Errorf("Could not connect to database: %v", err)
	}

	// Create route with GET method for test
	userStore := NewStorage(db)
	userController := NewController(userStore)
	Routes(app, userController)

	// // add health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	// // Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest("GET", test.route, nil)

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, 1)

		assert.Equal(test.expectedCode, resp.StatusCode)

	}
	// assert.Equal(t, 123, 123, "they should be equal")

}
