package meditation

import (
	"cmd/http/main.go/internal/user"
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

type createMeditationRequest struct {
	MeditationTime int   `json:"meditationTime" bson:"meditationTime"`
	EndTime        int64 `json:"endTime" bson:"endTime"`
}

type createMeditationResponse struct {
	ID string `json:"id"`
}

type getAllMeditationResponse []struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	UserID         string             `json:"userId" bson:"userId"`
	MeditationTime int                `json:"meditationTime" bson:"meditationTime"`
	EndTime        int64              `json:"endTime" bson:"endTime"`
}

type getMeditationResponse struct {
	Meditations []MeditationDB `json:"meditations"`
}

// @Summary Create meditation.
// @Description Creates a new meditation.
// @Tags meditation
// @Accept */*
// @Produce json
// @Param meditation body createMeditationRequest true "Meditation to create"
// @Param userId header string false "User ID"
// @Success 200 {object} createMeditationResponse
// @Router /meditation [post]
func (t *Controller) create(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req createMeditationRequest

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

	//TODO correct error handling
	// Create meditation record
	id, err := t.storage.Create(req, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to Create Meditation",
			"err":     err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(createMeditationResponse{
		ID: id,
	})
}

// @Summary Get meditation sessions
// @Description Fetch one or multiple meditation sessions.
// @Tags meditation
// @Param id query string false "Meditation ID"
// @Param startTime query int64 false "start time"
// @Param endTime query int64 false "end time"
// @Param durationStart query int64 false "duration start time"
// @Param durationEnd query int64 false "duration end time"
// @Param userId header string false "User ID"
// @Produce json
// @Success 200 {object} getMeditationResponse
// @Router /meditation [Get]
func (t *Controller) get(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	//parse Query values
	meditationId := c.Query("id")
	//map for time parameters
	times := map[string]int64{
		"startTime":     convertToInt64(c.Query("startTime")),
		"endTime":       convertToInt64(c.Query("endTime")),
		"startDuration": convertToInt64(c.Query("durationStart")),
		"durationEnd":   convertToInt64(c.Query("durationEnd")),
	}
	if meditationId != "" {
		// Get particular meditation
		meditation, err := t.storage.Get(meditationId, c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get meditation",
			})
		}
		// convert to array
		return c.Status(fiber.StatusOK).JSON(
			[]MeditationDB{meditation},
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
		// all meditations for a user between a time range and duration
		meditations, err := t.storage.GetAllOfOneUserBetweenTimeAndDuration(userId, times, c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get meditations in time range",
				"err":     err,
			})
		}
		if len(meditations) != 0 {
			//check if user is allowed to get this meditation
			if meditations[0].UserID != userId {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "User is not allowed to get this meditation",
				})
			}
			return c.Status(fiber.StatusOK).JSON(meditations)
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
