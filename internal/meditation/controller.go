package meditation

import (
	"cmd/http/main.go/internal/user"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	UserID         string `json:"userId" bson:"userId"`
	MeditationTime string `json:"meditationTime" bson:"meditationTime"`
	EndTime        string `json:"endTime" bson:"endTime"`
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
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	MeditationTime int                `json:"meditationTime" bson:"meditationTime"`
	EndTime        int64              `json:"endTime" bson:"endTime"`
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

// @Summary Get a meditation session
// @Description fetch a single meditation session.
// @Tags meditation
// @Param id path string true "Meditation ID"
// @Param userId header string true "User ID"
// @Produce json
// @Success 200 {object} getMeditationResponse
// @Router /meditation/{id} [Get]
func (t *Controller) get(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")

	meditationID := c.Params("meditationID")
	if meditationID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to Get meditations",
		})
	}

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	// Get meditation record
	meditation, err := t.storage.Get(meditationID, c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to fetch meditation",
		})
	}

	//check if user is allowed to get this meditation
	if meditation.UserID != userId {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User is not allowed to get this meditation",
		})
	}

	return c.Status(fiber.StatusOK).JSON(meditation)
}

// @Summary Get all meditation session
// @Description fetch all meditation sessions of a user.
// @Tags meditation
// @Param userId header string false "User ID"
// @Produce json
// @Success 200 {object} getAllMeditationResponse
// @Router /meditation/getAll [Get]
func (t *Controller) getAll(c *fiber.Ctx) error {

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

	// Get all meditations of a user
	meditations, err := t.storage.GetAllOfOneUser(userId, c.Context())
	if err != nil {
		fmt.Println("err", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to fetch meditation",
		})
	}
	return c.Status(fiber.StatusOK).JSON(meditations)
}
