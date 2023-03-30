package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	storage *Storage
}

func NewController(storage *Storage) *Controller {
	return &Controller{
		storage: storage,
	}
}

type createUserRequest struct {
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
	Email       string `json:"email" bson:"email"`
	ID          string `json:"id" bson:"id"`
}

type createUserResponse struct {
	ID string `json:"id"`
}

type getUserRequest struct {
	ID string `json:"id"`
}

// @Summary Create one user.
// @Description creates one user.
// @Tags users
// @Accept */*
// @Produce json
// @Param user body createUserRequest true "User to create"
// @Success 200 {object} createUserResponse
// @Router /users [post]
func (t *Controller) create(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req createUserRequest

	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"err":     err,
		})
	}

	//create user
	id, err := t.storage.create(req, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(createUserResponse{
		ID: id,
	})
}

// @Summary Get a user.
// @Description fetch a user by id.
// @Tags users
// @Accept */*
// @Produce json
// @Success 200 {object} userDB
// @Router /users [get]
func (t *Controller) get(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")

	var req getUserRequest

	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"err":     err,
		})
	}

	// get users
	user, err := t.storage.get(req.ID, c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get users",
		})
	}

	return c.JSON(user)
}

// @Summary Get all users.
// @Description fetch every user available.
// @Tags users
// @Accept */*
// @Produce json
// @Success 200 {object} []userDB
// @Router /users [get]
func (t *Controller) getAll(c *fiber.Ctx) error {
	// get all users
	users, err := t.storage.getAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get users",
		})
	}

	return c.JSON(users)
}
