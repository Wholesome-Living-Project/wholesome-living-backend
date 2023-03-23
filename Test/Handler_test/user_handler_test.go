package handlers_test

import (
	"bytes"
	"encoding/json"
	_ "encoding/json"
	"github.com/Wholesome-Living-Project/wholesome-living-backend/database"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	_ "net/http/httptest"
	"testing"

	_ "github.com/stretchr/testify/assert"

	"github.com/Wholesome-Living-Project/wholesome-living-backend/handlers"
)

func TestHandleCreateUser(t *testing.T) {
	// set up test app
	app := fiber.New()
	database.StartMongoDB()
	defer database.CloseMongoDB()

	// create request body
	requestBody := handlers.CreateUserDTO{
		ID:          "1",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1990-01-01",
		Email:       "john.doe@example.com",
	}

	// create request
	req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(requestBody.String))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// create response recorder
	rec := httptest.NewRecorder()

	// perform request
	app.Post("/user", handlers.HandleCreateUser)
	app.Handler()(req)

	// check response status code
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v want %v", status, http.StatusOK)
	}

	// parse response body
	var res handlers.CreateUserResDTO
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("could not parse response body: %v", err)
	}
}
