package elevator

import (
	"bytes"
	"cmd/http/main.go/internal/progress"
	"cmd/http/main.go/internal/storage"
	"cmd/http/main.go/internal/user"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	app        *fiber.App
	store      *Storage
	userStore  *user.Storage
	testUserId string
	elevatorId string
}

func (suite *Suite) SetupSuite() {
	// Define Fiber app.
	app := fiber.New()
	MONGODB_URI := "mongodb://localhost:27017"
	MONGODB_NAME := "testing"

	db, err := storage.BootstrapMongo(MONGODB_URI, MONGODB_NAME, 10*time.Second)

	if err != nil {
		suite.T().Errorf("Could not connect to database: %v", err)
	}

	if err := db.Collection("users").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

	userStore := user.NewStorage(db)
	suite.userStore = userStore
	progressStore := progress.NewStorage(db)

	suite.store = NewStorage(db)
	elevatorController := NewController(suite.store, userStore, progressStore)
	Routes(app, elevatorController)

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

	_, err := suite.store.db.Collection("users").Indexes().DropOne(context.Background(), "testId")
	if err != nil {
		log.Println("Error: ", err)
	}

	if err := suite.store.db.Collection("elevator").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

	// create a test user (just for userId purposes)
	testId, err := suite.userStore.Create(user.CreateUserRequest{
		ID:        "testId",
		FirstName: "test",
		LastName:  "testId",
	}, context.Background())

	if err != nil {
		suite.T().Errorf("Could not create test user: %v", err)
		testId = "testId"
	}

	suite.testUserId = testId

	// create test evelevators
	elevatorId, err := suite.store.Create(CreateElevatorRequest{
		Stairs:       true,
		AmountStairs: 12,
		HeightGain:   12,
	}, "testId", context.Background())

	if err != nil {
		suite.T().Errorf("Could not create test elevator: %v", err)
	}

	suite.elevatorId = elevatorId
}

func (suite *Suite) TestPost() {
	route := "/elevator"

	type Body struct {
		Stairs       bool    `json:"stairs"`
		AmountStairs int     `json:"amountStairs"`
		HeightGain   float64 `json:"heightGain"`
	}

	tests := []struct {
		missingHeader bool
		userId        string
		description   string
		expectedCode  int
		body          interface{}
	}{
		{
			description:  "Create successfully",
			body:         Body{true, 12, 12},
			expectedCode: fiber.StatusCreated,
		},
		{
			description:  "Another sucess test",
			body:         Body{true, 100, 12},
			expectedCode: fiber.StatusCreated,
		},
		{
			description:  "User does not exist",
			userId:       "doesntexist",
			body:         Body{true, 100, 12},
			expectedCode: fiber.StatusNotFound,
		},
		{
			description:   "Missing userId header",
			missingHeader: true,
			body:          Body{true, 100, 12},
			expectedCode:  fiber.StatusBadRequest,
		},
		{
			description:  "Wrongly setting amount of stairs",
			body:         Body{false, 12, 12},
			expectedCode: fiber.StatusNotFound,
		},
	}

	for _, test := range tests {
		bodyJson, err := json.Marshal(test.body)
		if err != nil {
			suite.T().Errorf("Could not marshal user: %v", err)
		}

		req := httptest.NewRequest("POST", route, bytes.NewReader(bodyJson))

		if test.userId != "" {
			req.Header.Set("userId", test.userId)
		} else {
			req.Header.Set("userId", suite.testUserId)
		}

		if test.missingHeader {
			req.Header.Del("userId")
		}

		resp, err := suite.app.Test(req, -1)
		if err != nil {
			suite.T().Errorf("Could not make request: %v", err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		success := suite.Equal(test.expectedCode, resp.StatusCode, "Error for (%v): %v ", test.description, string(body[:]))
		if !success {
			suite.T().Fail()
		}
		suite.T().Logf("[✔] (%v) passed", test.description)
	}
}

func (suite *Suite) TestGet() {
	route := "/elevator"

	tests := []struct {
		missingHeader bool
		userId        string
		description   string
		expectedCode  int
		query         map[string]string
	}{
		{
			description: "simple test one",
			query: map[string]string{
				"id": suite.elevatorId,
			},
			expectedCode: fiber.StatusOK,
		},
		{
			description: "empty id",
			query: map[string]string{
				"id": "",
			},
			expectedCode: fiber.StatusOK,
		},
		{
			description:   "Missing userId header",
			missingHeader: true,
			query:         map[string]string{},
			expectedCode:  fiber.StatusBadRequest,
		},
		{
			description:  "id does not exist",
			query:        map[string]string{"id": "nonexistingid"},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			description:  "User does not exist",
			userId:       "doesntexist",
			query:        map[string]string{},
			expectedCode: fiber.StatusNotFound,
		},
		{
			description: "Invalid body",
			query: map[string]string{
				"id":            "",
				"startTime":     "kasldkfjasldf",
				"endTime":       "",
				"startDuration": "aaldkfjaslkdfjaskldj",
				"durationEnd":   "",
				"minGain":       "",
				"maxGain":       "",
			},
			expectedCode: fiber.StatusOK,
		},
		{
			description: "Valid body",
			query: map[string]string{
				"startTime": "10",
			},
			expectedCode: fiber.StatusOK,
		},
	}

	for _, test := range tests {
		url := url.URL{
			Path: route,
		}

		// Add query
		q := url.Query()
		for key, value := range test.query {
			q.Add(key, string(value))
		}
		url.RawQuery = q.Encode()

		req := httptest.NewRequest("GET", url.String(), nil)

		// Add header
		if test.userId != "" {
			req.Header.Set("userId", test.userId)
		} else {
			req.Header.Set("userId", suite.testUserId)
		}

		if test.missingHeader {
			req.Header.Del("userId")
		}

		// test request
		resp, err := suite.app.Test(req, -1)
		if err != nil {
			suite.T().Errorf("Could not make request: %v", err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		success := suite.Equal(test.expectedCode, resp.StatusCode, "Error for (%v): %v ", test.description, string(body[:]))
		if !success {
			suite.T().Fail()
		}
		suite.T().Logf("[✔] (%v) passed", test.description)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTripTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
