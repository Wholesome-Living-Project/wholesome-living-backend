package handlers

import (
	"github.com/Wholesome-Living-Project/wholesome-living-backend/database"
	"github.com/gofiber/fiber/v2"
)

// @Summary Create medi.
// @Description Creates a new meditation.
// @Tags meditation
// @Accept json
// @Param userId, time, CreateMeditationDTO true "Todo to create"
// @Produce json
// @Success 200
// @Router /meditation [post]
func HandleCreateMeditation(c *fiber.Ctx) error {
	// get the user from the request body
	newMeditation := new(CreateMeditationDTO)

	// validate the request body
	if err := c.BodyParser(newMeditation); err != nil {
		return c.Status(400).JSON(fiber.Map{"bad input": err.Error()})
	}

	// insert the user into the database
	coll := database.GetCollection("meditation")
	res, err := coll.InsertOne(c.Context(), newMeditation)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"internal server error": err.Error()})
	}

	// return the inserted user
	return c.Status(200).JSON(fiber.Map{"inserted_id": res.InsertedID})
}

type CreateMeditationDTO struct {
	ID          string `json:"id" bson:"_id"`
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
	Email       string `json:"email" bson:"email"`
	CreatedAt   string `json:"createdAt" bson:"createdAt"`
}
