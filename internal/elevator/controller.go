package elevator

import (
	"cmd/http/main.go/internal/progress"
	"cmd/http/main.go/internal/settings"
	"cmd/http/main.go/internal/user"
	"github.com/gofiber/fiber/v2"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
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

type createElevatorRequest struct {
	Stairs       bool `json:"stairs" bson:"stairs"`
	AmountStairs int  `json:"amountStairs" bson:"amountStairs"`
}

type createElevatorResponse struct {
	ID string `json:"id"`
}

type getElevatorResponse struct {
	Meditations []ElevatorDB `json:"elevators"`
}

// @Summary Create elevator.
// @Description Creates a new elevator.
// @Tags elevator
// @Accept */*
// @Produce json
// @Param elevator body createElevatorRequest true "Elevator to create"
// @Param userId header string true "User ID"
// @Success 200 {object} createElevatorResponse
// @Router /elevator [post]
func (t *Controller) create(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req createElevatorRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"err":     err,
		})
	}

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	//check if user exists
	_, err := t.userStorage.Get(userId, c.Context())

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User does not exist",
			"err":     err,
		})
	}

	id, err := t.storage.Create(req, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to Create Elevator",
			"err":     err.Error(),
		})
	}
	err = t.progressStorage.AddExperience(userId, c.Context(), settings.PluginNameElevator, float64(req.AmountStairs/10))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(createElevatorResponse{
		ID: id,
	})
}

// @Summary Get elevator sessions
// @Description Fetch one or multiple elevator sessions.
// @Tags elevator
// @Param id query string false "Elevator ID"
// @Param startTime query int64 false "start time"
// @Param endTime query int64 false "end time"
// @Param durationStart query int64 false "duration start time"
// @Param durationEnd query int64 false "duration end time"
// @Param userId header string false "User ID"
// @Produce json
// @Success 200 {object} getElevatorResponse
// @Router /elevator [Get]
func (t *Controller) get(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	//parse Query values
	elevatorId := c.Query("id")
	//map for time parameters
	times := map[string]int64{
		"startTime":     convertToInt64(c.Query("startTime")),
		"endTime":       convertToInt64(c.Query("endTime")),
		"startDuration": convertToInt64(c.Query("durationStart")),
		"durationEnd":   convertToInt64(c.Query("durationEnd")),
	}
	if elevatorId != "" {
		// Get particular elevator
		elevator, err := t.storage.Get(elevatorId, c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get elevator",
			})
		}
		// convert to array
		return c.Status(fiber.StatusOK).JSON(
			[]ElevatorDB{elevator},
		)
	}

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}
	if userId != "" {
		_, err := t.userStorage.Get(userId, c.Context())

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User does not exist",
				"err":     err,
			})
		}
		// all elevators items for a user between a time range and duration
		elevators, err := t.storage.GetAllOfOneUserBetweenTimeAndDuration(userId, times, c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get elevators in time range",
				"err":     err,
			})
		}
		if len(elevators) != 0 {
			//check if user is allowed to get this elevator
			if elevators[0].UserID != userId {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "User is not allowed to get this elevator",
				})
			}
			return c.Status(fiber.StatusOK).JSON(elevators)
		}

	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Query constraints do not yield any results"})
}

func convertToInt64(value string) int64 {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		// Handle the error here if necessary
		// For example, you can assign a default value or log the error
		return 0
	}
	return intValue
}
