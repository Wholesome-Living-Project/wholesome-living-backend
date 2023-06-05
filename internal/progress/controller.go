package progress

import (
	"cmd/http/main.go/internal/user"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	storage     *Storage
	userStorage *user.Storage
}

func NewController(storage *Storage, userStorage *user.Storage) *Controller {
	return &Controller{
		storage:     storage,
		userStorage: userStorage,
	}
}

// @Summary Get progress nad level for a user.
// @Description fetch progress and level for a user.
// @Tags progress
// @param userId header string true "User ID"
// @Produce json
// @Success 200
// @Router /progress [get]
func (t *Controller) get(c *fiber.Ctx) error {
	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}
	// Get plugin from query

	settings, err := t.storage.Get(userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get settings, because: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(settings)
}
