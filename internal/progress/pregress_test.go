package progress

import (
	"cmd/http/main.go/internal/storage"
	"cmd/http/main.go/internal/user"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http/httptest"
	"testing"
	"time"
)

type ProgressSuite struct {
	suite.Suite
	app         *fiber.App
	storage     *Storage
	userStorage *user.Storage
	testUserId  string
}

func (suite *ProgressSuite) SetupSuite() {
	app := fiber.New()
	MONGODB_URI := "mongodb://localhost:27017"
	MONGODB_NAME := "testing"

	db, err := storage.BootstrapMongo(MONGODB_URI, MONGODB_NAME, 10*time.Second)
	if err != nil {
		suite.T().Errorf("Could not connect to database: %v", err)
	}

	suite.storage = NewStorage(db)
	suite.userStorage = user.NewStorage(db) // Replace with your own initialization
	ProgressController := NewController(suite.storage, suite.userStorage)
	app.Get("/progress", ProgressController.get)

	suite.app = app
	progressUser, err := suite.userStorage.Create(user.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		ID:        "progressUser",
	}, context.Background())
	suite.testUserId = progressUser

	if err != nil {
		suite.T().Errorf("Could not create test user: %v", err)
	}

	log.Println("SETUP DONE")
}

func (suite *ProgressSuite) TestInvalidGet() {
	req := httptest.NewRequest("GET", "/progress", nil)
	req.Header.Set("userId", "invalid")
	resp, _ := suite.app.Test(req)

	suite.Equal(404, resp.StatusCode, "Should return HTTP 400")
}

func TestProgressSuite(t *testing.T) {
	suite.Run(t, new(ProgressSuite))
}

// Tear down the database after all tests are done
func (suite *ProgressSuite) TearDownSuite() {
	// Drop the user collection
	if err := suite.storage.db.Collection("users").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

	// Drop the progress collection
	if err := suite.storage.db.Collection("progress").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}
}
