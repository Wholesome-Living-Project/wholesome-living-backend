package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	storage *UserStorage
}

func NewUserController(storage *UserStorage) *UserController {
	return &UserController{
		storage: storage,
	}
}

type createUserRequest struct {
	Name string `json:"name"`
}

type createUserResponse struct {
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
func (t *UserController) create(c *fiber.Ctx) error {
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
	id, err := t.storage.createUser(req.Name, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(createUserResponse{
		ID: id,
	})
}

// @Summary Get all users.
// @Description fetch every user available.
// @Tags users
// @Accept */*
// @Produce json
// @Success 200 {object} []userDB
// @Router /users [get]
func (t *UserController) getAll(c *fiber.Ctx) error {
	// get all users
	users, err := t.storage.getAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get users",
		})
	}

	return c.JSON(users)
}
