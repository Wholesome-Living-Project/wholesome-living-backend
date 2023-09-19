package settings

import (
	"cmd/http/main.go/internal/user"
	"fmt"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	"github.com/robfig/cron"
	"reflect"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	storage     *Storage
	userStorage *user.Storage
	cron        *cron.Cron
}

// an onboarding request can create settings for all plugins
type CreateSettingsRequest struct {
	EnabledPlugins []PluginName       `json:"enabledPlugins" bson:"enabledPlugins"`
	Meditation     MeditationSettings `json:"meditation" bson:"meditation"`
	Finance        FinanceSettings    `json:"finance" bson:"finance"`
	Elevator       ElevatorSettings   `json:"elevator" bson:"elevator"`
}

func NewController(storage *Storage, userStorage *user.Storage, cron *cron.Cron) *Controller {
	return &Controller{
		storage:     storage,
		userStorage: userStorage,
		cron:        cron,
	}
}

// an onboarding request can create settings for all plugins
type CreateOnboardingSettingResponse struct {
	EnabledPlugins []PluginName       `json:"enabledPlugins" bson:"enabledPlugins"`
	Meditation     MeditationSettings `json:"meditation" bson:"meditation"`
	Finance        FinanceSettings    `json:"finance" bson:"finance"`
	Elevator       ElevatorSettings   `json:"elevator" bson:"elevator"`
}

type createOnboardingResponse struct {
	ID string `json:"id"`
}

type getPluginSettingResponse struct {
	SettingsDB
}

// @Summary Create onboarding in backend, set settings.
// @Description Creates settings for a user.
// @Tags settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body CreateSettingsRequest true "onboarding to create"
// @Success 200 {object} createOnboardingResponse
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

	return c.Status(fiber.StatusCreated).JSON(createOnboardingResponse{http})
}

// @Summary Get plugin settings for a user.
// @Description fetch plugin settings for a user.
// @Tags settings
// @param userId header string true "User ID"
// @Param plugin query string false "Plugin name"
// @Produce json
// @Success 200 {object} getPluginSettingResponse
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
		f := reflect.Indirect(r).FieldByName(cases.Title(language.English).String(plugin))

		if f.IsValid() {
			return c.Status(fiber.StatusOK).JSON(f.Interface())
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Plugin does not exist",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(getPluginSettingResponse{settings})
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
	c.Request().Header.Set("Content-Type", "application/json")

	userId := string(c.Request().Header.Peek("userId"))
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing userId header",
		})
	}

	user, err := t.userStorage.Get(userId, c.Context())

	// check if user has token
	if user.ExpoPushToken == "" {
		return nil
	}

	// To check the token is valid
	pushToken, err := expo.NewExponentPushToken(user.ExpoPushToken)
	if err != nil {
		panic(err)
	}

	var req MeditationSettings
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"err":     err,
		})
	}

	err = t.AddMeditationNotificationInterval(req, pushToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not schedule meditation push notification: " + err.Error(),
		})
	}
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
	return c.Status(fiber.StatusOK).JSON(http)
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
		if err.Error() == "User not found" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "User not found" + err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete settings, because: " + err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

// Notify @Summary push notification
// @Description send a push notification to a user's device.
// @Tags users
// @Accept */*
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (t *Controller) Notify(token expo.ExponentPushToken, message string, title string) error {
	// Create a new Expo SDK client
	client := expo.NewPushClient(nil)

	// Publish message
	response, err := client.Publish(
		&expo.PushMessage{
			To:       []expo.ExponentPushToken{token},
			Body:     message,
			Data:     map[string]string{"withSome": "data"},
			Sound:    "default",
			Title:    title,
			Priority: expo.DefaultPriority,
		},
	)

	// Check errors
	if err != nil {
		panic(err)
	}

	// Validate responses
	if response.ValidateResponse() != nil {
		fmt.Println(response.PushMessage.To, "failed")
	}

	return nil
}

func (t *Controller) AddMeditationNotificationInterval(req MeditationSettings, token expo.ExponentPushToken) error {

	if req.PeriodNotifications == "Day" {
		t.Notify(token, "Did you meditate today?", "It is time meditate!")
	}

	return nil
}
