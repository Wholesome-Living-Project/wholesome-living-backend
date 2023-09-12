package settings

import (
	"cmd/http/main.go/internal/user"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	storage     *Storage
	userStorage *user.Storage
}

// an onboarding request can create settings for all plugins
type CreateSettingsRequest struct {
	EnabledPlugins []PluginName       `json:"enabledPlugins" bson:"enabledPlugins"`
	Meditation     MeditationSettings `json:"meditation" bson:"meditation"`
	Finance        FinanceSettings    `json:"finance" bson:"finance"`
	Elevator       ElevatorSettings   `json:"elevator" bson:"elevator"`
}

func NewController(storage *Storage, userStorage *user.Storage) *Controller {
	return &Controller{
		storage:     storage,
		userStorage: userStorage,
	}
}

// an onboarding request can create settings for all plugins
type CreateOnboardingSettingResponse struct {
	EnabledPlugins []PluginName       `json:"enabledPlugins" bson:"enabledPlugins"`
	Meditation     MeditationSettings `json:"meditation" bson:"meditation"`
	Finance        FinanceSettings    `json:"finance" bson:"finance"`
	Elevator       ElevatorSettings   `json:"elevator" bson:"elevator"`
}

// TODO for each Plugin a creat endpoint

/*
type createInvestmentResponse struct {
	ID string `json:"id"`
}

type getPluginSettingResponse struct {
	EnabledPlugins []PluginName       `json:"enabledPlugins" bson:"enabledPlugins"`
	Meditation     MeditationSettings `json:"meditation" bson:"meditation"`
	Finance        FinanceSettings    `json:"finance" bson:"finance"`
	Elevator       ElevatorSettings   `json:"elevator" bson:"elevator"`
}
*/

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
			"message": "Could not create onboarding, because: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(http)
}

// @Summary Get plugin settings for a user.
// @Description fetch plugin settings for a user.
// @Tags settings
// @param userId header string true "User ID"
// @Param plugin query string false "Plugin name"
// @Produce json
// @Success 200 {object} getPluginSettingsResponse
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

	settings, err := t.storage.Get(userId, plugin, c.Context())
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
// @Description Creates settings for a user for the finance Plugin.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body FinanceSettings true "onboarding to create"
// @Success 201
// @Router /settings/finance [post]
func (t *Controller) createFinanceSettings(c *fiber.Ctx) error {
	return t.createPluginSettings(&FinanceSettings{}, c)
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
	return t.createPluginSettings(&ElevatorSettings{}, c)
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
	return t.createPluginSettings(&MeditationSettings{}, c)
}

func (t *Controller) createPluginSettings(settingType SingleSetting, c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	if err := c.BodyParser(settingType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body" + err.Error(),
		})
	}

	err := t.storage.CreatePluginSettings(settingType, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create " + settingType.getName() + " settings: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON("Created")
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
	return t.updatePluginSettings(&FinanceSettings{}, c)
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
	return t.updatePluginSettings(&MeditationSettings{}, c)
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
	return t.updatePluginSettings(&ElevatorSettings{}, c)
}

func (t *Controller) updatePluginSettings(settingType SingleSetting, c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")

	// Check if the user is logged in
	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	if err := c.BodyParser(settingType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	http, err := t.storage.UpdatePluginSettings(settingType, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"message": "Could not update" +
					settingType.getName() +
					" settings: " + err.Error(),
			},
		)
	}
	return c.Status(fiber.StatusCreated).JSON(http)
}

// @Summary Delete plugin-settings of a user.
// @Description Delete plugin-settings for a user if plugin is "" delete all settings.
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
	err := t.storage.Delete(userId, plugin, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete settings because: " + err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
