package user

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"testing"
	"time"

	"cmd/http/main.go/internal/storage"

	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	app        *fiber.App
	store      *Storage
	testUserId string
}

func (suite *Suite) SetupSuite() {
	// Define Fiber app.
	app := fiber.New()
	MONGODB_URI := "mongodb://localhost:27017"
	MONGODB_NAME := "testing-user"

	db, err := storage.BootstrapMongo(MONGODB_URI, MONGODB_NAME, 10*time.Second)

	if err != nil {
		suite.T().Errorf("Could not connect to database: %v", err)
	}

	if err := db.Collection("users").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

	suite.store = NewStorage(db)
	userController := NewController(suite.store)
	Routes(app, userController)

	// // add health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	suite.app = app

	log.Println("SETUP DONE")

}

func (suite *Suite) BeforeTest(suiteName, testName string) {
	if err := suite.store.db.Collection("users").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

	// create a test user (just for userId purposes)
	testId, err := suite.store.Create(CreateUserRequest{
		FirstName: "test",
		ID:        "testId",
	}, context.Background())

	if err != nil {
		suite.T().Errorf("Could not create test user: %v", err)
	}

	suite.testUserId = testId

}

func (suite *Suite) TestGetFast() {
	tests := []struct {
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
	}{
		{
			description:  "Health check",
			route:        "/health",
			expectedCode: 200,
		},
		{
			description:  "get all users (fast)",
			route:        "/users",
			expectedCode: 200,
		},
		{
			description:  "get nonex-users (fast)",
			route:        "/users/123",
			expectedCode: 404,
		},
		{
			description:  "get existing users (fast)",
			route:        "/users/" + suite.testUserId,
			expectedCode: 200,
		},
	}

	for _, test := range tests {
		suite.T().Log(test.description)
		req := httptest.NewRequest("GET", test.route, nil)
		resp, _ := suite.app.Test(req, 1)
		suite.Equal(test.expectedCode, resp.StatusCode)
	}
}

func (suite *Suite) TestPost() {
	route := "/users"

	tests := []struct {
		description  string
		expectedCode int
		user         map[string]string
	}{
		{
			description: "Create successfully",
			user: map[string]string{
				"userName":   "test",
				"nonkeyword": "body",
				"id":         "123",
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			description: "ID already exists",
			user: map[string]string{
				"username": "test",
				"id":       "123",
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			description: "ID is empty",
			user: map[string]string{
				"username": "test",
				"id":       "",
			},
			expectedCode: fiber.StatusCreated, // TODO whut?
		},
	}

	for _, test := range tests {
		bodyJson, err := json.Marshal(test.user)
		if err != nil {
			suite.T().Errorf("Could not marshal user: %v", err)
		}

		req := httptest.NewRequest("POST", route, bytes.NewReader(bodyJson))

		resp, err := suite.app.Test(req, -1)
		if err != nil {
			suite.T().Errorf("Could not make request: %v", err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		suite.Equal(test.expectedCode, resp.StatusCode, "Error for Create: %v", string(body[:]))
	}
}

func (suite *Suite) TestPutUser() {
	route := "/users"

	tests := []struct {
		description  string
		userId       string
		expectedCode int
		user         map[string]string
	}{
		{
			description: "User changed successfully",
			userId:      suite.testUserId, // will set no header
			user: map[string]string{
				"userName": "test changed",
			},
			expectedCode: fiber.StatusOK,
		},
		{
			description: "Update-user without userId",
			userId:      "", // will set no header
			user: map[string]string{
				"username": "test will not change",
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			description: "Update non existing user",
			userId:      "nonexistant", // will set no header
			user: map[string]string{
				"username": "test will not change",
			},
			expectedCode: fiber.StatusNotFound,
		},
		{
			description: "ID already exists", // TODO remove this abiltity
			userId:      suite.testUserId,
			user: map[string]string{
				"username": "test",
				"id":       suite.testUserId,
			},
			expectedCode: fiber.StatusOK,
		},
		{
			description:  "Empty body",
			userId:       suite.testUserId,
			user:         map[string]string{},
			expectedCode: fiber.StatusOK, //
		},
	}

	for _, test := range tests {
		suite.T().Log(test.description)
		bodyJson, err := json.Marshal(test.user)
		if err != nil {
			suite.T().Errorf("Could not marshal user: %v", err)
		}

		req := httptest.NewRequest("PUT", route, bytes.NewReader(bodyJson))

		if test.userId != "" {
			req.Header.Set("userId", test.userId)
		}

		resp, err := suite.app.Test(req, -1)
		if err != nil {
			suite.T().Errorf("Could not make request: %v", err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		suite.Equal(test.expectedCode, resp.StatusCode, "%v", string(body[:]))
	}

}

func (suite *Suite) TestDeleteUser() {
	tests := []struct {
		description  string
		userId       string
		expectedCode int
	}{
		{
			description:  "Delete existing user",
			userId:       suite.testUserId,
			expectedCode: fiber.StatusOK,
		},
		{
			description:  "Delete non-existing user",
			userId:       "nonexistent",
			expectedCode: fiber.StatusNotFound,
		},
	}

	for _, test := range tests {
		suite.T().Log(test.description)
		req := httptest.NewRequest("DELETE", "/users/"+test.userId, nil)
		resp, err := suite.app.Test(req, -1)
		if err != nil {
			suite.T().Errorf("Could not make request: %v", err)
		}

		suite.Equal(test.expectedCode, resp.StatusCode)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// cleanup after the the suite
func (suite *Suite) TearDownSuite() {
	if err := suite.store.db.Collection("users").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}
}
