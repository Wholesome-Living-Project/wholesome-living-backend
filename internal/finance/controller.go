package finance

import (
	"cmd/http/main.go/internal/progress"
	"cmd/http/main.go/internal/settings"
	"cmd/http/main.go/internal/user"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	storage         *Storage
	userStorage     *user.Storage
	progressStorage *progress.Storage
}

func NewController(storage *Storage, userStorage *user.Storage, progressStorage *progress.Storage) *Controller {
	return &Controller{
		storage:         storage,
		userStorage:     userStorage,
		progressStorage: progressStorage,
	}
}

type CreateSpendingRequest struct {
	Amount       float64 `json:"amount" bson:"amount"`
	Saving       float64 `json:"saving" bson:"saving"`
	SpendingTime int64   `json:"spendingTime" bson:"spendingTime"`
	Description  string  `json:"description" bson:"description"`
}

type createSpendingResponse struct {
	ID string `json:"id"`
}

type getInvestmentResponse struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       string             `json:"userId" bson:"userId"`
	SpendingTime int64              `json:"spendingTime" bson:"spendingTime"`
	Amount       float64            `json:"amount" bson:"amount"`
	Saving       float64            `json:"saving" bson:"saving"`
	Description  string             `json:"description" bson:"description"`
}

// @Summary Create a spending.
// @Description Creates a new spending.
// @Tags finance
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param investment body createSpendingRequest true "spending to create"
// @Success 200 {object} createSpendingResponse
// @Router /finance [post]
func (t *Controller) create(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req CreateSpendingRequest
	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"err":     err,
		})
	}

	//TODO correct error handling
	id, err := t.storage.create(req, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to create",
			"err":     err,
		})
	}

	if err != nil {
		return err
	}
	err = t.progressStorage.AddExperience(userId, c.Context(), settings.PluginNameFinance, float64(req.Saving)/2)
	return c.Status(fiber.StatusCreated).JSON(createSpendingResponse{
		ID: id,
	})
}

// @Summary Query Investments with the user ID, start time and end time.
// @Description Query Investments with the user ID, start time and end time.
// @Tags finance
// @param userId header string true "User ID"
// @Param id query string false "investment ID"
// @Param startTime query int64 false "start time"
// @Param endTime query int64 false "end time"
// @Produce json
// @Success 200 {object} getInvestmentResponse
// @Router /finance [get]
func (t *Controller) get(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	//parse Query values
	particularInvestment := c.Query("id")
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")
	var startTime, endTime int64
	var err error

	if startTimeStr != "" {
		startTime, err = strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid startTime parameter",
				"err":     err,
			})
		}
	}

	if endTimeStr != "" {
		endTime, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid endTime parameter",
				"err":     err,
			})
		}
	}

	if particularInvestment != "" {
		// Get particular investment investment
		investment, err := t.storage.get(particularInvestment, c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get investment",
			})
		}
		// Convert FinanceDb to getInvestmentResponse
		investmentResponse := getInvestmentResponse{
			ID:           investment.ID,
			UserID:       investment.UserID,
			SpendingTime: investment.SpendingTime,
			Amount:       investment.Amount,
			Saving:       investment.Saving,
			Description:  investment.Description,
		}
		return c.JSON([]getInvestmentResponse{investmentResponse})
	}

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	} else {
		_, err := t.userStorage.Get(userId, c.Context())

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User does not exist",
				"err":     err,
			})
		}
		if startTimeStr == "" && endTimeStr == "" {
			// all investments for a user
			investments, err := t.storage.getAllOfOneUser(userId, c.Context())

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to get investments",
					"err":     err,
				})

			}
			return c.Status(fiber.StatusOK).JSON(investments)
		}
		if startTimeStr != "" || endTimeStr != "" {
			// all investments for a user between a time range
			// Todo if startTime is given and endTime is not given, then return all investments after startTime
			investments, err := t.storage.getAllOfOneUserBetweenTime(userId, startTime, endTime, c.Context())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to get investments in time range",
					"err":     err,
				})
			}
			return c.Status(fiber.StatusOK).JSON(investments)
		}
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "Invalid request body",
	})
}
