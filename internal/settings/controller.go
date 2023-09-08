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
type CreateSettingsRequest struct {
	// A list with the Plugins that the user has enabled.
	EnabledPlugins []PluginName `json:"enabledPlugins" bson:"enabledPlugins"`
	// The user's settings for the meditation plugin.
	Meditation MeditationSettings `json:"meditation" bson:"meditation"`
	// The user's settings for the finance plugin.
	Finance  FinanceSettings  `json:"finance" bson:"finance"`
	Elevator ElevatorSettings `json:"elevator" bson:"elevator"`
}

// TODO for each Plugin a creat endpoint

type createInvestmentResponse struct {
	ID string `json:"id"`
}

type getSettingsResponse struct {
	// A list with the Plugins that the user has enabled.
	EnabledPlugins []PluginName `json:"enabledPlugins" bson:"enabledPlugins"`
	// The user's settings for the meditation plugin.
	Meditation MeditationSettings `json:"meditation" bson:"meditation"`
	// The user's settings for the finance plugin.
	Finance  FinanceSettings  `json:"finance" bson:"finance"`
	Elevator ElevatorSettings `json:"elevator" bson:"elevator"`
}

// @Summary Create onboarding in backend, set settings.
// @Description Creates settings for a user.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body CreateSettingsRequest true "onboarding to create"
// @Success 200 {object} createInvestmentResponse
// @Router /settings [post]
func (t *Controller) createOnboarding(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req CreateSettingsRequest
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
// @Success 200 {object} getSettingsResponse
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
		if err.Error() == "mongo: no documents in result" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Plugin does not exist",
			})
		}
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
// @Description Creates settings for a user for the finance Plugin.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body FinanceSettings true "onboarding to create"
// @Success 201
// @Router /settings/finance [post]
func (t *Controller) createFinanceSettings(c *fiber.Ctx) error {
	return t.createPluginSettings(c, "finance")
}

// @Summary Create settings for the elevator plugin.
// @Description Creates settings for a user for the elevator Plugin.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body ElevatorSettings true "onboarding to create"
// @Success 201
// @Router /settings/elevator [post]
func (t *Controller) createElevatorSettings(c *fiber.Ctx) error {
	return t.createPluginSettings(c, "elevator")
}

// @Summary Create settings for the meditation Plugin.
// @Description Creates settings for the meditation plugin
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body MeditationSettings true "onboarding to create"
// @Success 201
// @Router /settings/meditation [post]
func (t *Controller) createMeditationSettings(c *fiber.Ctx) error {
	return t.createPluginSettings(c, "meditation")
}

func (t *Controller) createPluginSettings(c *fiber.Ctx, pluginName string) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req interface{}
	switch pluginName {
	case "finance":
		var financeReq FinanceSettings
		req = &financeReq
	case "meditation":
		var meditationReq MeditationSettings
		req = &meditationReq
	case "elevator":
		var elevatorReq ElevatorSettings
		req = &elevatorReq
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid plugin name",
		})
	}

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	http, err := t.storage.CreatePluginSettings(req, userId, pluginName, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create " + pluginName + " settings: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(http)
}

// updateFinanceSettings
// @Summary updateFinanceSettings settings for the finance plugin.
// @Description Update settings for a user for onr Plugin.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body FinanceSettings true "onboarding to create"
// @Success 200
// @Router /settings/finance [put]
func (t *Controller) updateFinanceSettings(c *fiber.Ctx) error {
	return t.updatePluginSettings(c, "finance")
}

// @Summary Update settings for the meditation Plugin.
// @Description Update settings for a user
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body MeditationSettings true "onboarding to create"
// @Success 200
// @Router /settings/meditation [put]
func (t *Controller) updateMeditationSettings(c *fiber.Ctx) error {
	return t.updatePluginSettings(c, "meditation")
}

// @Summary Update settings for the elevator Plugin.
// @Description Update settings for the elevator Plugin.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body ElevatorSettings true "onboarding to create"
// @Success 200
// @Router /settings/elevator [put]
func (t *Controller) updateElevatorSettings(c *fiber.Ctx) error {
	return t.updatePluginSettings(c, "elevator")
}

func (t *Controller) updatePluginSettings(c *fiber.Ctx, pluginName string) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req interface{}
	switch pluginName {
	case "finance":
		var financeReq FinanceSettings
		req = &financeReq
	case "meditation":
		var meditationReq MeditationSettings
		req = &meditationReq
	case "elevator":
		var elevatorReq ElevatorSettings
		req = &elevatorReq
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid plugin name",
		})
	}

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	http, err := t.storage.UpdatePluginSettings(req, userId, pluginName, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create " + pluginName + " settings: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(http)
}

// @Summary Delete settings of a user.
// @Description Delete settings for a user.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param plugin query string false "Plugin name"
// @Success 201
// @Router /settings [delete]
func (t *Controller) delete(c *fiber.Ctx) error {
	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}
	plugin := c.Query("plugin")
	_, err := t.storage.Delete(userId, c.Context(), plugin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete settings, because: " + err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
