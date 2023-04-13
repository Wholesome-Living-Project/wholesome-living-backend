package finance

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

type createSpendingRequest struct {
	UserID       string `json:"userId" bson:"userId"`
	Amount       string `json:"amount" bson:"amount"`
	SpendingTime string `json:"endTime" bson:"endTime"`
}

type createSpendingResponse struct {
	ID string `json:"id"`
}

type getSpendingResponse struct {
	UserID       primitive.ObjectID `json:"userId" bson:"userId"`
	SpendingTime string             `json:"spendingTime" bson:"spendingTime"`
	Amount       string             `json:"amount" bson:"amount"`
}

// @Summary Create spending.
// @Description Creates a new spending.
// @Tags finance
// @Accept */*
// @Produce json
// @Param spending body createSpendingRequest true "spending to create"
// @Success 200 {object} createSpendingResponse
// @Router /spending [post]
func (t *Controller) create(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req createSpendingRequest

	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"err":     err,
		})
	}

	//TODO correct error handling
	// create spending record
	id, err := t.storage.create(req, c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to create spending",
			"err":     err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(createSpendingResponse{
		ID: id,
	})
}

// @Summary Get a spending session
// @Description fetch a single spending session.
// @Tags finance
// @Param id path string true "spending ID"
// @Produce json
// @Success 200 {object} getSpendingResponse
// @Router /spending/{id} [get]
func (t *Controller) get(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")

	spendingID := c.Params("spendingID")
	if spendingID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to get spendings",
		})
	}

	// create spending record
	user, err := t.storage.get(spendingID, c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to fetch spending",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary Get all spending session
// @Description fetch all spending's of one user.
// @Tags finance
// @Param userID path string true "User ID"
// @Produce json
// @Success 200 {object} getSpendingResponse
// @Router /spending/getAll/{userID} [get]
func (t *Controller) getAll(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Provide an ID",
		})
	}
	//TODO correct error handling
	// create spending record

	// get all spendings of a user
	spendings, err := t.storage.getAllOfOneUser(userID, c.Context())
	if err != nil {
		fmt.Println("errrr", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to fetch spending",
		})
	}
	return c.Status(fiber.StatusOK).JSON(spendings)
}
