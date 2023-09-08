package settings

import (
	"bytes"
	"cmd/http/main.go/internal/storage"
	"cmd/http/main.go/internal/user"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http/httptest"
	"testing"
	"time"
)

type SettingsSuite struct {
	suite.Suite
	app         *fiber.App
	store       *Storage      // Initialize as per your setup
	userStorage *user.Storage // Initialize as per your setup
	testUserId  string
	testUserId2 string
}

func (suite *SettingsSuite) SetupSuite() {
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

	suite.store = NewStorage(db)
	suite.userStorage = user.NewStorage(db) // Replace with your own initialization
	SettingsController := NewController(suite.store, suite.userStorage)
	Routes(app, SettingsController)

	// // add health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	suite.app = app

	log.Println("SETUP DONE")
}

func (suite *SettingsSuite) BeforeTest(suiteName, testName string) {
	if err := suite.store.db.Collection("users").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

	// create a test user (just for userId purposes)
	settingsUser, err := suite.userStorage.Create(user.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		ID:        "settingsUser",
	}, context.Background())

	if err != nil {
		suite.T().Errorf("Could not create test user: %v", err)
	}

	suite.testUserId = settingsUser

}

func (suite *SettingsSuite) TestCreateOnboardingGetAndDelete() {
	reqBody := CreateSettingsRequest{
		EnabledPlugins: []PluginName{"meditation", "elevator"},
		Meditation: MeditationSettings{
			MeditationTimeGoal:  10,
			Notifications:       true,
			AmountNotifications: 3,
			PeriodNotifications: "Day",
		},
		Finance: FinanceSettings{
			Notifications:       true,
			AmountNotifications: 0,
			PeriodNotifications: "Day",
			Strategy:            "Round",
			StrategyAmount:      0,
			InvestmentGoal:      0,
			InvestmentTimeGoal:  0,
		},
		Elevator: ElevatorSettings{
			Notifications:       true,
			AmountNotifications: 3,
			PeriodNotifications: "Day",
			Goal:                100,
		}}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/settings", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("userId", suite.testUserId)

	resp, _ := suite.app.Test(req)

	suite.Equal(201, resp.StatusCode, "\nStatus::"+resp.Status+"\n", "Should return HTTP 201")

	req2 := httptest.NewRequest("GET", "/settings", nil)
	req2.Header.Set("userId", suite.testUserId)

	resp2, _ := suite.app.Test(req2)

	suite.Equal(200, resp2.StatusCode, resp2.Status, "Should return HTTP 200")
	req3 := httptest.NewRequest("DELETE", "/settings", nil)
	req3.Header.Set("userId", suite.testUserId)

	resp3, _ := suite.app.Test(req3)

	suite.Equal(200, resp3.StatusCode, "\nStatus::"+resp3.Status+"\n", "Should return HTTP 200")

}

func (suite *SettingsSuite) TestCreateFinanceSettingsAndDelete() {
	reqBody := FinanceSettings{
		Notifications:       true,
		AmountNotifications: 3,
		PeriodNotifications: "Day",
		Strategy:            "Round",
		StrategyAmount:      5,
		InvestmentGoal:      100,
		InvestmentTimeGoal:  10,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/settings/finance", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("userId", suite.testUserId)

	resp, _ := suite.app.Test(req)
	suite.Equal(201, resp.StatusCode, "\nStatus::"+resp.Status+"\n", "Should return HTTP 201")

	//test put
	req2 := httptest.NewRequest("PUT", "/settings/finance", bytes.NewReader(reqBodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("userId", suite.testUserId)
	resp2, _ := suite.app.Test(req2)
	suite.Equal(201, resp2.StatusCode, "\nStatus::"+resp2.Status+"\n", "Should return HTTP 200")

	req3 := httptest.NewRequest("DELETE", "/settings", nil)
	req3.Header.Set("userId", suite.testUserId)

	resp3, _ := suite.app.Test(req3)

	suite.Equal(200, resp3.StatusCode, "\nStatus::"+resp3.Status+"\n", "Should return HTTP 200")

}

func (suite *SettingsSuite) TestGetInvalidPlugin() {
	req := httptest.NewRequest("GET", "/settings", nil)
	req.Header.Set("userId", suite.testUserId)
	req.Header.Set("plugin", "invalidPlugin")

	resp, _ := suite.app.Test(req)

	suite.Equal(400, resp.StatusCode, "\nStatus::"+resp.Status+"\n", "Should return HTTP 400")

}
func TestSettingsSuite(t *testing.T) {
	suite.Run(t, new(SettingsSuite))
}

// tear down the database after all tests are done
func (suite *SettingsSuite) TearDownSuite() {
	// user
	if err := suite.store.db.Collection("users").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}
	// settings
	if err := suite.store.db.Collection("settings").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

}
