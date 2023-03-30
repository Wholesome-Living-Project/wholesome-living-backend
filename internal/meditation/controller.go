package meditation

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	storage *Storage
}

func NewController(storage *Storage) *Controller {
	return &Controller{
		storage: storage,
	}
}

type createMeditationRequest struct {
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	MeditationTime string             `json:"meditationTime" bson:"meditationTime"`
	EndTime        string             `json:"endTime" bson:"endTime"`
}

type createMeditationResponse struct {
	ID string `json:"id"`
}

type meditationResponse struct {
	MeditationID primitive.ObjectID `json:"meditationId" bson:"meditationId"`
}

// @Summary Create meditation.
// @Description Creates a new meditation.
// @Tags meditation
// @Accept */*
// @Produce json
// @Param meditation body createMeditationRequest true "Meditation to create"
// @Success 200 {object} createMeditationResponse
// @Router /meditation [post]

func (t *Controller) create(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req createMeditationRequest

	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"err":     err,
		})
	}

	// create meditation record
	id, err := t.storage.create(req, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
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
// @Produce json
// @Success 200 {object} meditationResponse
// @Router /meditation/{id} [get]
func (t *Controller) get(c *fiber.Ctx) error {
	// get the id from the request params
	var err error = nil
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON("{}")
	}
	return c.Status(fiber.StatusOK).JSON("{}")
}
