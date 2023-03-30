package meditation

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Controller struct {
	storage *MediationStorage
}

func NewController(storage *MediationStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

type createMeditationRequest struct {
	MeditationTime string `json:"meditationTime" bson:"meditationTime"`
}

type createMeditationResponse struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}

type meditationResponse struct {
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	LastName       string             `json:"completed" bson:"completed"`
	CreatedAt      string             `json:"date" bson:"date"`
	MeditationTime string             `json:"meditationTime" bson:"meditationTime"`
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

	// TODO use new storage functions for geting and seting db info

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

// @Summary Get a meditation session
// @Description fetch a single meditation session.
// @Tags meditation
// @Param id path string true "Meditation ID"
// @Produce json
// @Success 200 {object} meditationResponse
// @Router /meditation/{id} [get]
func (t *Controller) get(c *fiber.Ctx) error {
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
