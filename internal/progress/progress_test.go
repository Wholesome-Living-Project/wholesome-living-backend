package progress

import (
	"cmd/http/main.go/internal/storage"
	"cmd/http/main.go/internal/user"
	"context"
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
}

func (suite *Suite) SetupSuite() {
	// Define Fiber app.
	app := fiber.New()
	MONGODB_URI := "mongodb://localhost:27017"
	MONGODB_NAME := "testing-progress"

	db, err := storage.BootstrapMongo(MONGODB_URI, MONGODB_NAME, 10*time.Second)

	if err != nil {
		suite.T().Errorf("Could not connect to database: %v", err)
	}

	if err := db.Collection("users").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

	userStore := user.NewStorage(db)
	suite.userStore = userStore

	suite.store = NewStorage(db)
	progCont := NewController(suite.store, userStore)
	Routes(app, progCont)

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

	if err := suite.store.db.Collection("progress").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

	// create a test user (just for userId purposes)
	testId := "testId"
	_, err = suite.userStore.Get(testId, context.Background())

	if err != nil {
		_, err := suite.userStore.Create(user.CreateUserRequest{
			ID:        "testId",
			FirstName: "test",
			LastName:  "testId",
		}, context.Background())

		if err != nil {
			suite.T().Errorf("Could not create test user: %v", err)
		}

	}

	suite.testUserId = testId
}

func (suite *Suite) TestGet() {
	route := "/progress"

	tests := []struct {
		missingHeader bool
		userId        string
		description   string
		expectedCode  int
		query         map[string]string
	}{
		{
			description:   "Missing userId header",
			missingHeader: true,
			query:         map[string]string{},
			expectedCode:  fiber.StatusBadRequest,
		},
		{
			description:  "User does not exist",
			userId:       "doesntexist",
			query:        map[string]string{},
			expectedCode: fiber.StatusNotFound,
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
func TestTripTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
