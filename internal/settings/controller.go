package settings

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

// an onboarding request can create settings for all plugins
type createOnboardingRequest struct {
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
}

// @Summary Create onboarding in backend, set settings.
// @Description Creates settings for a user.
// @Tags Settings
// @Accept */*
// @Produce json
// @param userId header string true "User ID"
// @Param settings body createOnboardingRequest true "onboarding to create"
// @Success 200 {object} createInvestmentResponse
// @Router /settings [post]
func (t *Controller) createOnboarding(c *fiber.Ctx) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req createOnboardingRequest
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

	// TODO check if user exists

	// TODO check if user already has onboarding settings
	http, err := t.storage.CreateOnboarding(req, userId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create onboarding, because " + http,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(http)

}

// @Summary Get a single investment
// @Description fetch a single investment session.
// @Tags finance
// @param userId header string true "User ID"
// @Param id query string false "investment ID"
// @Param startTime query int64 false "start time"
// @Param endTime query int64 false "end time"
// @Produce json
// @Success 200 {object} getInvestmentResponse
// @Router /investment [get]
func (t *Controller) get(c *fiber.Ctx) error {

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "Invalid request body",
	})
}
