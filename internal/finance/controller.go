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

type createInvestmentRequest struct {
	UserID         string `json:"userId" bson:"userId"`
	Amount         string `json:"amount" bson:"amount"`
	InvestmentTime string `json:"endTime" bson:"endTime"`
}

type createInvestmentResponse struct {
	ID string `json:"id"`
}

type getInvestmentResponse struct {
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	InvestmentTime string             `json:"investmentTime" bson:"investmentTime"`
	Amount         string             `json:"amount" bson:"amount"`
}

// @Summary Create a investment.
// @Description Creates a new investment.
// @Tags finance
// @Accept */*
// @Produce json
// @Param investment body createInvestmentRequest true "investment to create"
// @Success 200 {object} createInvestmentResponse
// @Router /investment [post]
func (t *Controller) create(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req createInvestmentRequest

	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"err":     err,
		})
	}

	//TODO correct error handling
	// create investment record
	id, err := t.storage.create(req, c.Context())
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

// @Summary Get a single investment
// @Description fetch a single investment session.
// @Tags finance
// @Param id path string true "investment ID"
// @Produce json
// @Success 200 {object} getInvestmentResponse
// @Router /investment/{id} [get]
func (t *Controller) get(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")

	investmentID := c.Params("investmentID")
	if investmentID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to get investments",
		})
	}

	// create investment record
	user, err := t.storage.get(investmentID, c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Failed to fetch investment",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary Get all investments of a user
// @Description fetch all investment's of one user.
// @Tags finance
// @Param userID path string true "User ID"
// @Produce json
// @Success 200 {object} getInvestmentResponse
// @Router /investment/getAll/{userID} [get]
func (t *Controller) getAll(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Provide an ID",
		})
	}
	//TODO correct error handling
	// create investment record

	// get all investments of a user
	investments, err := t.storage.getAllOfOneUser(userID, c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to fetch investment",
		})
	}
	return c.Status(fiber.StatusOK).JSON(investments)
}
