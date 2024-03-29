package meditation

import (
	"cmd/http/main.go/internal/progress"
	"cmd/http/main.go/internal/settings"
	"cmd/http/main.go/internal/user"
	"strconv"

	"github.com/gofiber/fiber/v2"
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

type CreateMeditationRequest struct {
	MeditationTime int   `json:"meditationTime" bson:"meditationTime"`
	EndTime        int64 `json:"endTime" bson:"endTime"`
}

type createMeditationResponse struct {
	ID string `json:"id"`
}

// TODO remove if not needed
/*
type getAllMeditationResponse []struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	UserID         string             `json:"userId" bson:"userId"`
	MeditationTime int                `json:"meditationTime" bson:"meditationTime"`
	EndTime        int64              `json:"endTime" bson:"endTime"`
}
*/

// @Summary Create meditation.
// @Description Creates a new meditation.
// @Tags meditation
// @Accept */*
// @Produce json
// @Param meditation body CreateMeditationRequest true "Meditation to create"
// @Param userId header string true "User ID"
// @Success 200 {object} createMeditationResponse
// @Router /meditation [post]
func (t *Controller) create(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req CreateMeditationRequest

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
	err = t.progressStorage.AddExperience(userId, c.Context(), settings.PluginNameMeditation, float64(req.MeditationTime))
	if err != nil {
		return err
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
// @Success 200 {object} []MeditationDB
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
	} else {
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

		return c.Status(fiber.StatusOK).JSON(meditations)
	}
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
