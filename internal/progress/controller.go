package progress

import (
	"cmd/http/main.go/internal/user"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
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
// @Param plugin query string false "Plugin name"
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
	plugin := c.Query("plugin")

	settings, err := t.storage.Get(userId, c.Context(), plugin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get settings, because: " + err.Error(),
		})
	}

	if plugin != "" {
		r := reflect.ValueOf(settings)
		f := reflect.Indirect(r).FieldByName(strings.Title(plugin))

		if f.IsValid() {
			return c.Status(fiber.StatusOK).JSON(f.Interface())
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid plugin name",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(settings)
}
