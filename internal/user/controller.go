package user

import (
	"fmt"
)

type PluginName string

const (
	PluginNameMeditation PluginName = "meditation"
	PluginNameFinance    PluginName = "finance"
	PluginNameElevator   PluginName = "elevator"
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
	ID          string `json:"id" bson:"_id"`
}

type createUserResponse struct {
	ID string `json:"id"`
}

type getUserRequest struct {
	ID string `json:"id"`
}

type updateUserRequest struct {
	FirstName   string       `json:"firstName" bson:"firstName"`
	LastName    string       `json:"lastName" bson:"lastName"`
	DateOfBirth string       `json:"dateOfBirth" bson:"dateOfBirth"`
	Email       string       `json:"email" bson:"email"`
	Plugins     []PluginName `json:"plugins" bson:"plugins"`
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

	//Create user
	_, err := t.storage.Create(req, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to Create user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(createUserResponse{
		ID: req.ID,
	})
}

// @Summary Get a user.
// @Description fetch a user by id.
// @Tags users
// @Accept */*
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserDB
// @Router /users/{id} [Get]
func (t *Controller) get(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get users
	user, err := t.storage.Get(id, c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to Get users",
		})
	}

	return c.JSON(user)
}

// @Summary Get all users.
// @Description fetch every user available.
// @Tags users
// @Accept */*
// @Produce json
// @Success 200 {object} []UserDB
// @Router /users [Get]
func (t *Controller) getAll(c *fiber.Ctx) error {
	// Get all users
	users, err := t.storage.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to Get users",
		})
	}

	return c.JSON(users)
}

// @Summary Update a user.
// @Description update a user by id.
// @Tags users
// @Accept */*
// @Produce json
// @Param user body updateUserRequest true "User to update"
// @Param userId header string false "User ID"
// @Success 200 {object} UserDB
// @Router /users [put]
func (t *Controller) update(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	// Parse the update request from the request body
	var req updateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"err":     err,
		})
	}

	// Fetch the existing user from the database
	user, err := t.storage.Get(userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Does Not Exist",
		})
	}

	// Update the user object with the new values
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.DateOfBirth != "" {
		user.DateOfBirth = req.DateOfBirth
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Plugins != nil {
		// check if plugin is valid
		for _, plugin := range req.Plugins {
			if plugin != PluginNameMeditation && plugin != PluginNameFinance && plugin != PluginNameElevator {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Invalid plugin name",
				})
			}
		}
		user.Plugins = req.Plugins
	}

	// Update the user in the database
	result, err := t.storage.Update(user, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user",
		})
	}

	return c.JSON(result)
}
