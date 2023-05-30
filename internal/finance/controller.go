package finance

import (
	"cmd/http/main.go/internal/user"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

type Controller struct {
	storage     *Storage
	userStorage *user.Storage
}

func NewController(storage *Storage, userStorage *user.Storage) *Controller {
	return &Controller{
		storage:     storage,
		userStorage: userStorage,
	}
}

type createInvestmentRequest struct {
	Amount         int    `json:"amount" bson:"amount"`
	InvestmentTime int64  `json:"investmentTime" bson:"investmentTime"`
	Description    string `json:"description" bson:"description"`
}

type createInvestmentResponse struct {
	ID string `json:"id"`
}

type getInvestmentResponse struct {
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	InvestmentTime int64              `json:"investmentTime" bson:"investmentTime"`
	Amount         int                `json:"amount" bson:"amount"`
}

// @Summary Create a investment.
// @Description Creates a new investment.
// @Tags finance
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param investment body createInvestmentRequest true "investment to create"
// @Success 200 {object} createInvestmentResponse
// @Router /finance [post]
func (t *Controller) create(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req createInvestmentRequest
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
	// create investment record
	id, err := t.storage.create(req, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to create",
			"err":     err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(createInvestmentResponse{
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
		return c.JSON(investment)
	}

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	} //check if user exists
	if userId != "" {
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
