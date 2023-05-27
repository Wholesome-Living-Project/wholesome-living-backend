package settings

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type MeditationSettings struct {
	// The user's meditation time goal.
	MeditationTimeGoal  int              `json:"meditationTime" bson:"meditationTime"`
	Notifications       bool             `json:"notifications" bson:"notifications"`
	AmountNotifications int              `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications NotificationType `json:"periodNotifications" bson:"periodNotifications"`
}
type NotificationType string

const (
	NotificationTypeDay   NotificationType = "Day"
	NotificationTypeMonth NotificationType = "Month"
	NotificationTypeWeek  NotificationType = "Week"
)

type StrategyType string

const (
	StrategyTypeRound   StrategyType = "Round"
	StrategyTypePlus    StrategyType = "Plus"
	StrategyTypePercent StrategyType = "Percent"
)

type FinanceSettings struct {
	Notifications       bool             `json:"notifications" bson:"notifications"`
	AmountNotifications int              `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications NotificationType `json:"periodNotifications" bson:"periodNotifications"`
	Strategy            StrategyType     `json:"strategy" bson:"strategy"`
	StrategyAmount      int              `json:"strategyAmount" bson:"strategyAmount"`
	// The user's investment goal.
	InvestmentGoal int `json:"investmentGoal" bson:"investmentGoal"`
	// The user's investment time goal.
	InvestmentTimeGoal int `json:"investmentTimeGoal" bson:"investmentTimeGoal"`
	// Interest rate of investment TODO maybe there can be multiple different interest rates for different investmens
}

type SettingsDB struct {
	// A list with the Plugins that the user has enabled.
	EnabledPlugins []pluginType `json:"enabledPlugins" bson:"enabledPlugins"`
	// The user's settings for the meditation plugin.
	Meditation MeditationSettings `json:"meditation" bson:"meditation"`
	// The user's settings for the finance plugin.
	Finance FinanceSettings `json:"finance" bson:"finance"`
	ID      string          `json:"id" bson:"_id"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Get(userId string, ctx context.Context, plugin string) (SettingsDB, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")
	settingsRecord := SettingsDB{}

	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := userResult.Err(); err != nil {
		return settingsRecord, errors.New("User not found!")
	}

	if plugin != "" {
		// check if plugin exists
		if isValidPlugins(plugin) == false {
			return settingsRecord, errors.New("Plugin not found!")
		}
		// Get certain plugin info
		cursor := collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": plugin})
		if err := cursor.Err(); err != nil {
			return settingsRecord, err
		}
		if err := cursor.Decode(&settingsRecord); err != nil {
			return settingsRecord, err
		}
	}
	cursor := collection.FindOne(ctx, bson.M{"_id": userId})
	if err := cursor.Err(); err != nil {
		return settingsRecord, err
	}

	if err := cursor.Decode(&settingsRecord); err != nil {
		return settingsRecord, err
	}
	return settingsRecord, nil
}

// TODO: Settings should be when onboarding is made or a user choses a ned plugin

func (s *Storage) CreateOnboarding(request createSettingsRequest, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")

	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := userResult.Err(); err != nil {
		return "", errors.New("User not found!")
	}
	// Check if user already has onboarding settings
	userSettings := collection.FindOne(ctx, bson.M{"_id": userId})
	if userSettings.Err() == nil {
		return "", errors.New("User already has onboarding settings")
	}
	if err := userResult.Err(); err != nil {
		return "User not found", err
	}
	if validateSettings(request) != nil {
		return "Invalid settings", validateSettings(request)
	}
	settings := SettingsDB{
		ID:             userId,
		EnabledPlugins: request.EnabledPlugins,
		Meditation:     request.Meditation,
		Finance:        request.Finance,
	}

	result, err := collection.InsertOne(ctx, settings)

	if err != nil {
		return result.InsertedID.(string), err
	}
	return "Inserted", err
}
func (s *Storage) createPluginSettings(request interface{}, userId string, pluginName string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")

	// Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := userResult.Err(); err != nil {
		return "", errors.New("User not found!")
	}

	// Check if user already has onboarding settings
	settings := collection.FindOne(ctx, bson.M{"_id": userId})
	if settings.Err() == nil {
		// Check if user already has the specified plugin settings
		if collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": pluginName}).Err() == nil {
			return "", errors.New("User already has " + pluginName + " settings, use put")
		} else {
			// Create plugin settings and keep the other settings
			var settingsRecord SettingsDB
			if err := settings.Decode(&settingsRecord); err != nil {
				return "", err
			}
			if validateSettings(request) != nil {
				return "Invalid settings", errors.New("Invalid settings: " + validateSettings(request).Error())
			}
			newSettings := SettingsDB{
				ID:             userId,
				EnabledPlugins: append(settingsRecord.EnabledPlugins, pluginType(pluginName)),
			}

			switch pluginName {
			case "finance":
				newSettings.Finance = *(request.(*FinanceSettings))
				newSettings.Meditation = settingsRecord.Meditation
			case "meditation":
				newSettings.Meditation = *(request.(*MeditationSettings))
				newSettings.Finance = settingsRecord.Finance
			}

			// Update the settings with the new settings
			result := collection.FindOneAndUpdate(ctx, bson.M{"_id": userId}, bson.M{"$set": newSettings})
			if result.Err() != nil {
				return "", result.Err()
			} else {
				return "Updated", nil
			}
		}

	} else {
		// Create new settings since the user has no settings yet
		if validateSettings(request) != nil {
			return "Invalid settings", errors.New("Invalid settings: " + validateSettings(request).Error())
		}
		settings := SettingsDB{
			ID:             userId,
			EnabledPlugins: []pluginType{pluginType(pluginName)},
		}

		switch pluginName {
		case "finance":
			settings.Finance = *(request.(*FinanceSettings))
		case "meditation":
			settings.Meditation = *(request.(*MeditationSettings))
		}

		// Insert the settings
		_, err := collection.InsertOne(ctx, settings)
		if err != nil {
			return "", err
		}
		return "Inserted", err
	}

	return "Something went wrong", nil
}

func (s *Storage) updatePluginSettings(request interface{}, userId string, pluginName string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")

	// Check if user already has the specified plugin settings
	oldSettings := collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": pluginName})
	if oldSettings.Err() == nil {
		// Update plugin settings
		var settingsRecord SettingsDB
		if err := oldSettings.Decode(&settingsRecord); err != nil {
			return "", err
		}
		if validateSettings(request) != nil {
			return "Invalid settings", errors.New("Invalid settings: " + validateSettings(request).Error())
		}
		newSettings := SettingsDB{
			ID:             userId,
			EnabledPlugins: settingsRecord.EnabledPlugins,
		}
		switch pluginName {
		case "finance":
			newSettings.Finance = *(request.(*FinanceSettings))
			newSettings.Meditation = settingsRecord.Meditation
		case "meditation":
			newSettings.Meditation = *(request.(*MeditationSettings))
			newSettings.Finance = settingsRecord.Finance
		}
		if validateSettings(request) != nil {
			return "Invalid settings", validateSettings(request)
		}
		result := collection.FindOneAndUpdate(ctx, bson.M{"_id": userId}, bson.M{"$set": newSettings})
		if result.Err() != nil {
			return "", result.Err()
		} else {
			return "Updated", nil
		}
	} else {
		return "", errors.New("User does not have " + pluginName + " settings")
	}
}

// function to validate each setting of a plugin
func validateSettings(request interface{}) error {
	// Validate settings, can be extended
	newSettings := SettingsDB{}
	switch request.(type) {
	case *FinanceSettings:
		newSettings.Finance = *(request.(*FinanceSettings))
		if !isValidNotificationType(newSettings.Finance.PeriodNotifications) || !isValidStrategy(newSettings.Finance.Strategy) {
			return errors.New("Invalid notification type, either strat or notificcation")
		}

	case *MeditationSettings:
		newSettings.Meditation = *(request.(*MeditationSettings))
		if !isValidNotificationType(newSettings.Meditation.PeriodNotifications) {
			return errors.New("Invalid notification type")
		}
	}
	if !isValidPlugins(newSettings.EnabledPlugins) {
		return errors.New("Invalid plugin")
	}

	return nil
}
func isValidPlugins(plugin interface{}) bool {
	// check if plugins are valid
	// if type of plugin is string then it is only one plugin
	if reflect.TypeOf(plugin).Kind() == reflect.String {
		plugin := plugin.(string)
		if plugin != "finance" && plugin != "meditation" {
			return false
		}
		return true
	}
	plugins := plugin.([]pluginType)
	for _, plugin := range plugins {
		// TODO make this shit dynamic fmmmlll
		if plugin != "finance" && plugin != "meditation" {
			return false
		}

	}
	return true
}
func isValidNotificationType(notification NotificationType) bool {
	// Validate notification type
	switch notification {
	case NotificationTypeDay, NotificationTypeWeek, NotificationTypeMonth:
		return true
	default:
		return false
	}
}
func isValidStrategy(strat StrategyType) bool {
	// Validate notification type
	switch strat {
	case StrategyTypePercent, StrategyTypeRound, StrategyTypePlus:
		return true
	default:
		return false
	}
}
