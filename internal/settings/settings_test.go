package settings

import (
	"bytes"
	"cmd/http/main.go/internal/storage"
	"cmd/http/main.go/internal/user"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

type SettingsSuite struct {
	suite.Suite
	app         *fiber.App
	store       *Storage      // Initialize as per your setup
	userStorage *user.Storage // Initialize as per your setup
	testUserId  string
}

func (suite *SettingsSuite) SetupSuite() {
	app := fiber.New()
	MONGODB_URI := "mongodb://localhost:27017"
	MONGODB_NAME := "testing-setting"

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

	if err := suite.store.db.Collection("settings").Drop(context.Background()); err != nil {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	suite.Equal(500, resp.StatusCode, "Message: %v", string(body))

	//test put
	req2 := httptest.NewRequest("PUT", "/settings/finance", bytes.NewReader(reqBodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("userId", suite.testUserId)
	resp2, _ := suite.app.Test(req2)
	suite.Equal(200, resp2.StatusCode, "\nStatus::"+resp2.Status+"\n", "Should return HTTP 200")

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

	suite.Equal(500, resp.StatusCode, "\nStatus::"+resp.Status+"\n", "Should return HTTP 400")

}

// Test for missing userId header in createPluginSettings
func (suite *SettingsSuite) TestCreatePluginSettingsMissingUserId() {
	reqBody, _ := json.Marshal(&MeditationSettings{}) // Replace with actual data
	req := httptest.NewRequest("POST", "/settings/meditation", bytes.NewBuffer(reqBody))
	resp, _ := suite.app.Test(req)
	suite.Equal(400, resp.StatusCode)
}

// Test for invalid request body in createPluginSettings
func (suite *SettingsSuite) TestCreatePluginSettingsInvalidRequestBody() {
	req := httptest.NewRequest("POST", "/settings/meditation", bytes.NewBuffer([]byte("{")))
	req.Header.Add("userId", "someUserId")
	resp, _ := suite.app.Test(req)
	suite.Equal(400, resp.StatusCode)
}

// Test for missing userId header in delete settings
func (suite *SettingsSuite) TestDeleteSettingsMissingUserId() {
	req := httptest.NewRequest("DELETE", "/settings", nil)
	resp, _ := suite.app.Test(req)
	suite.Equal(400, resp.StatusCode)
}

func TestSettingsSuite(t *testing.T) {
	suite.Run(t, new(SettingsSuite))
}

func (suite *SettingsSuite) TestCreatePluginSettingsTwiceAndDeleteOneSetting() {
	reqBody := CreateSettingsRequest{
		EnabledPlugins: []PluginName{"elevator"},
		Meditation: MeditationSettings{
			MeditationTimeGoal:  0,
			Notifications:       true,
			AmountNotifications: 0,
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

	reqBody2 := MeditationSettings{
		MeditationTimeGoal:  10,
		Notifications:       true,
		AmountNotifications: 3,
		PeriodNotifications: "Day",
	}
	reqBodyBytes2, _ := json.Marshal(reqBody2)

	req2 := httptest.NewRequest("POST", "/settings/meditation", bytes.NewReader(reqBodyBytes2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("userId", suite.testUserId)

	resp2, _ := suite.app.Test(req2)
	suite.Equal(201, resp2.StatusCode, "\nStatus::"+resp2.Status+"\n", "Should return HTTP 201")

	req3 := httptest.NewRequest("POST", "/settings/meditation", bytes.NewReader(reqBodyBytes))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("userId", suite.testUserId)

	resp3, _ := suite.app.Test(req3)
	suite.Equal(500, resp3.StatusCode, "\nStatus::"+resp3.Status+"\n", "Should return HTTP 500")

	// Delete settings
	req5 := httptest.NewRequest("DELETE", "/settings?plugin=meditation", nil)
	req5.Header.Set("userId", suite.testUserId)
	resp5, _ := suite.app.Test(req5)

	suite.Equal(200, resp5.StatusCode, "\nStatus::"+resp5.Status+"\n", "Should return HTTP 200")

	req6 := httptest.NewRequest("DELETE", "/settings?plugin=elevator", nil)
	req6.Header.Set("userId", suite.testUserId)

	resp6, _ := suite.app.Test(req6)

	suite.Equal(200, resp6.StatusCode, "\nStatus::"+resp6.Status+"\n", "Should return HTTP 200")
}

func (suite *SettingsSuite) TestMissingUserId() {
	req := httptest.NewRequest("POST", "/settings/meditation", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := suite.app.Test(req)
	suite.Equal(fiber.StatusBadRequest, resp.StatusCode)

	req2 := httptest.NewRequest("POST", "/settings", nil)
	req2.Header.Set("Content-Type", "application/json")

	resp2, _ := suite.app.Test(req2)
	suite.Equal(fiber.StatusBadRequest, resp2.StatusCode)
}

func (suite *SettingsSuite) TestInvalidRequestBody() {
	reqBody := "{invalid_json"
	req := httptest.NewRequest("POST", "/settings/meditation", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("userId", "someUserId")

	resp, _ := suite.app.Test(req)
	suite.Equal(fiber.StatusBadRequest, resp.StatusCode)
}

func (suite *SettingsSuite) TestUpdateMeditationAndElevator() {
	// Workaround since AfterTest does not delete settings
	req6 := httptest.NewRequest("DELETE", "/settings", nil)
	req6.Header.Set("userId", suite.testUserId)

	resp6, _ := suite.app.Test(req6)

	body, err := io.ReadAll(resp6.Body)
	if err != nil {
		panic(err)
	}

	suite.Equal(500, resp6.StatusCode, "Message: %v", string(body))

	reqBody := CreateSettingsRequest{
		EnabledPlugins: []PluginName{"elevator", "meditation"},
		Meditation: MeditationSettings{
			MeditationTimeGoal:  24,
			Notifications:       true,
			AmountNotifications: 34,
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

	reqBody2 := MeditationSettings{
		MeditationTimeGoal:  10,
		Notifications:       true,
		AmountNotifications: 3,
		PeriodNotifications: "Day",
	}
	reqBodyBytes2, _ := json.Marshal(reqBody2)

	req2 := httptest.NewRequest("PUT", "/settings/meditation", bytes.NewReader(reqBodyBytes2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("userId", suite.testUserId)
	resp2, _ := suite.app.Test(req2)
	suite.Equal(200, resp2.StatusCode, "\nStatus::"+resp2.Status+"\n", "Should return HTTP 200")

	reqBody3 := ElevatorSettings{
		Notifications:       true,
		AmountNotifications: 3,
		PeriodNotifications: "Day",
		Goal:                100,
	}
	reqBodyBytes3, _ := json.Marshal(reqBody3)

	req3 := httptest.NewRequest("PUT", "/settings/elevator", bytes.NewReader(reqBodyBytes3))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("userId", suite.testUserId)
	resp3, _ := suite.app.Test(req3)
	suite.Equal(200, resp3.StatusCode, "\nStatus::"+resp3.Status+"\n", "Should return HTTP 200")

}

func (suite *SettingsSuite) TestBodyParserFail() {
	req := httptest.NewRequest("POST", "/settings/meditation", bytes.NewBufferString("{invalid_json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("userId", "someUserId")

	resp, _ := suite.app.Test(req)
	suite.Equal(fiber.StatusBadRequest, resp.StatusCode)
}

func (suite *SettingsSuite) TestGetSpecificPlugin() {
	reqBody := CreateSettingsRequest{
		EnabledPlugins: []PluginName{"elevator", "meditation", "finance"},
		Meditation: MeditationSettings{
			MeditationTimeGoal:  24,
			Notifications:       true,
			AmountNotifications: 34,
			PeriodNotifications: "Day",
		},
		Finance: FinanceSettings{
			Notifications:       true,
			AmountNotifications: 0,
			PeriodNotifications: "Day",
			Strategy:            "Round",
			StrategyAmount:      4,
			InvestmentGoal:      55,
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

	req2 := httptest.NewRequest("GET", "/settings?plugin=elevator", nil)
	req2.Header.Set("userId", suite.testUserId)
	resp2, _ := suite.app.Test(req2)
	suite.Equal(200, resp2.StatusCode, "\nStatus::"+resp2.Status+"\n", "Should return HTTP 200")

	req3 := httptest.NewRequest("GET", "/settings?plugin=meditation", nil)
	req3.Header.Set("userId", suite.testUserId)
	resp3, _ := suite.app.Test(req3)
	suite.Equal(200, resp3.StatusCode, "\nStatus::"+resp3.Status+"\n", "Should return HTTP 200")

	req4 := httptest.NewRequest("GET", "/settings?plugin=finance", nil)
	req4.Header.Set("userId", suite.testUserId)
	resp4, _ := suite.app.Test(req4)
	suite.Equal(200, resp4.StatusCode, "\nStatus::"+resp4.Status+"\n", "Should return HTTP 200")

	req5 := httptest.NewRequest("GET", "/settings?plugin=invalidPlugin", nil)
	req5.Header.Set("userId", suite.testUserId)
	resp5, _ := suite.app.Test(req5)
	suite.Equal(500, resp5.StatusCode, "\nStatus::"+resp5.Status+"\n", "Should return HTTP 400")

}

func (suite *SettingsSuite) TestPutSettingsMissingUserId() {
	req := httptest.NewRequest("PUT", "/settings/meditation", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := suite.app.Test(req)
	suite.Equal(fiber.StatusBadRequest, resp.StatusCode)
}

func (suite *SettingsSuite) TestDeleteSettingsUnavailableUserId() {
	req := httptest.NewRequest("DELETE", "/settings", nil)
	req.Header.Set("userId", "someUserId")

	resp, _ := suite.app.Test(req)
	suite.Equal(fiber.StatusInternalServerError, resp.StatusCode)
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
func (suite *SettingsSuite) AfterTest() {
	if err := suite.userStorage.Delete(suite.testUserId, context.Background()); err != nil {
		log.Println("Error: ", err)
	}
	// delete settings
	if err := suite.store.db.Collection("settings").Drop(context.Background()); err != nil {
		log.Println("Error: ", err)
	}

	if err := suite.store.Delete(suite.testUserId, "", context.Background()); err != nil {
		log.Println("Error: ", err)
	}
}
