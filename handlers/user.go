package handlers

import (
	"time"

	"github.com/Wholesome-Living-Project/wholesome-living-backend/database"
	"github.com/Wholesome-Living-Project/wholesome-living-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserDTO struct {
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
	Email       string `json:"email" bson:"email"`
	CreatedAt   string `json:"createdAt" bson:"createdAt"`
}

type CreateUserResDTO struct {
	InsertedId primitive.ObjectID `json:"inserted_id" bson:"_id"`
}

// @Summary Create a user.
// @Description creates a user and returns it.
// @Tags user
// @Accept json
// @Param user body CreateUserDTO true "Todo to create"
// @Produce json
// @Success 200 {object} CreateUserResDTO
// @Router /user [post]
func HandleCreateUser(c *fiber.Ctx) error {
	// get the user from the request body
	newUser := new(CreateUserDTO)

	// validate the request body
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"bad input": err.Error()})
	}

	newUser.CreatedAt = time.Now().GoString()

	// insert the user into the database
	coll := database.GetCollection("user")
	res, err := coll.InsertOne(c.Context(), newUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"internal server error": err.Error()})
	}

	// return the inserted user
	return c.Status(200).JSON(fiber.Map{"inserted_id": res.InsertedID})
}

type UserResDTO struct {
	ID          string `json:"id" bson:"_id"`
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
	Email       string `json:"email" bson:"email"`
	CreatedAt   string `json:"createdAt" bson:"createdAt"`
}

// @Summary Get a single user.
// @Description fetch a single user.
// @Tags user
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {object} models.User
// @Router /user/{id} [get]
func HandleGetUser(c *fiber.Ctx) error {
	// get the id from the request params
	id := c.Params("id")
	dbId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"invalid id": err.Error()})
	}

	// fetch the user from the database
	coll := database.GetCollection("users")
	filter := bson.M{"_id": dbId}
	var user models.User
	err = coll.FindOne(c.Context(), filter).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	// return the user
	return c.Status(200).JSON(user)
}
