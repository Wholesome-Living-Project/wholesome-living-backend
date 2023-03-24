package handlers

import (
	"github.com/Wholesome-Living-Project/wholesome-living-backend/database"
	"github.com/Wholesome-Living-Project/wholesome-living-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreateMeditationResDTO struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}
type MeditationResDTO struct {
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	LastName       string             `json:"completed" bson:"completed"`
	CreatedAt      string             `json:"date" bson:"date"`
	MeditationTime string             `json:"meditationTime" bson:"meditationTime"`
}
type CreateMeditationDTO struct {
	MeditationTime string `json:"meditationTime" bson:"meditationTime"`
}

// HandleCreateMeditation @Summary Create medi.
// @Description Creates a new meditation.
// @Tags meditation
// @Param meditation body CreateMeditationDTO true "Meditation to create"
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {object} models.Meditation
// @Router /meditation/{id} [post]
func HandleCreateMeditation(c *fiber.Ctx) error {
	newMeditation := new(MeditationResDTO)
	if err := c.BodyParser(newMeditation); err != nil {
		return c.Status(400).JSON(fiber.Map{"bad input": err.Error()})
	}

	givenUserId := c.Params("id")
	// decode the user id
	// todo not working yet check if it is a valid id
	dbId, err := primitive.ObjectIDFromHex(givenUserId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"invalid id": err.Error()})
	}
	// fetch the user from the database
	coll := database.GetCollection("user")
	filter := bson.M{"_id": dbId}
	var user models.User
	err = coll.FindOne(c.Context(), filter).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	newMeditation.CreatedAt = time.Now().GoString()
	newMeditation.UserID = dbId

	// insert the user into the database
	collMedi := database.GetCollection("meditation")
	res, err := collMedi.InsertOne(c.Context(), newMeditation)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"internal server error": err.Error()})
	}

	// return the inserted user
	return c.Status(200).JSON(fiber.Map{"inserted_id": res.InsertedID})
}

// HandleGetMeditation @Summary Get a meditation session
// @Description fetch a single meditation session.
// @Tags meditation
// @Param id path string true "Meditationo ID"
// @Produce json
// @Success 200 {object} models.User
// @Router /meditation/{id} [get]
func HandleGetMeditation(c *fiber.Ctx) error {
	// get the id from the request params
	id := c.Params("id")
	dbId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"invalid id": err.Error()})
	}

	// fetch the user from the database
	coll := database.GetCollection("meditation")
	filter := bson.M{"_id": dbId}
	var meditation models.Meditation
	err = coll.FindOne(c.Context(), filter).Decode(&meditation)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	// return the user
	return c.Status(200).JSON(meditation)
}
