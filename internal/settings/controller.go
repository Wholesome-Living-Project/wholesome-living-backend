package settings

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

// an onboarding request can create settings for all plugins
type createSettingsRequest struct {
	// A list with the Plugins that the user has enabled.
	EnabledPlugins []string `json:"enabledPlugins" bson:"enabledPlugins"`
	// The user's settings for the meditation plugin.
	Meditation MeditationSettings `json:"meditation" bson:"meditation"`
	// The user's settings for the finance plugin.
	Finance FinanceSettings `json:"finance" bson:"finance"`
}

// TODO for each Plugin a creat endpoint

type createInvestmentResponse struct {
	ID string `json:"id"`
}

type getInvestmentResponse struct {
	// A list with the Plugins that the user has enabled.
	EnabledPlugins []string `json:"enabledPlugins" bson:"enabledPlugins"`
	// The user's settings for the meditation plugin.
	Meditation MeditationSettings `json:"meditation" bson:"meditation"`
	// The user's settings for the finance plugin.
	Finance FinanceSettings `json:"finance" bson:"finance"`
}

// @Summary Create onboarding in backend, set settings.
// @Description Creates settings for a user.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body createSettingsRequest true "onboarding to create"
// @Success 200 {object} createInvestmentResponse
// @Router /settings [post]
func (t *Controller) createOnboarding(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req createSettingsRequest
	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	http, err := t.storage.CreateOnboarding(req, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create onboarding, because " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(http)

}

// @Summary Get settings for a user.
// @Description fetch settings for a user.
// @Tags settings
// @param userId header string true "User ID"
// @Param plugin query string false "Plugin name"
// @Produce json
// @Success 200 {object} getInvestmentResponse
// @Router /settings [get]
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

// @Summary Create settings for the finance plugin.
// @Description Creates settings for a user for onr Plugin.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body FinanceSettings true "onboarding to create"
// @Success 200 {object} createInvestmentResponse
// @Router /settings/finance [post]
func (t *Controller) createFinanceSettings(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req FinanceSettings
	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	http, err := t.storage.createFinanceSettings(req, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create finance settings, because " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(http)

}

// @Summary Create settings for the meditation Plugin.
// @Description Creates settings for a user
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body MeditationSettings true "onboarding to create"
// @Success 200 {object} createInvestmentResponse
// @Router /settings/meditation [post]
func (t *Controller) createMeditationSettings(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req MeditationSettings
	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	http, err := t.storage.createMeditationSettings(req, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create meditation settings, because " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(http)

}
