package settings

import (
	"context"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationType string

type PluginName string

// ADD NEW PLUGIN NAME HERE
const (
	PluginNameFinance    PluginName = "finance"
	PluginNameMeditation PluginName = "meditation"
	PluginNameElevator   PluginName = "elevator"
)

var validPlugins = map[PluginName]bool{
	PluginNameFinance:    true,
	PluginNameMeditation: true,
	PluginNameElevator:   true,
}

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

type MeditationSettings struct {
	// The user's meditation time goal.
	MeditationTimeGoal  int              `json:"meditationTimeGoal" bson:"meditationTimeGoal"`
	Notifications       bool             `json:"notifications" bson:"notifications"`
	AmountNotifications int              `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications NotificationType `json:"periodNotifications" bson:"periodNotifications"`
}

type ElevatorSettings struct {
	Notifications       bool             `json:"notifications" bson:"notifications"`
	AmountNotifications int              `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications NotificationType `json:"periodNotifications" bson:"periodNotifications"`
	Goal                int              `json:"goal" bson:"goal"`
}

type SettingsDB struct {
	// A list with the Plugins that the user has enabled.
	EnabledPlugins []PluginName `json:"enabledPlugins" bson:"enabledPlugins"`
	// The user's settings for the meditation plugin.
	Meditation MeditationSettings `json:"meditation" bson:"meditation"`
	// The user's settings for the finance plugin.
	Finance FinanceSettings `json:"finance" bson:"finance"`
	// The user's settings for the elevator plugin.
	Elevator ElevatorSettings `json:"elevator" bson:"elevator"`
	ID       string           `json:"id" bson:"_id"`
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
		pluginName := PluginName(plugin)
		if _, ok := validPlugins[pluginName]; !ok {
			return settingsRecord, errors.New("Plugin not found!")
		}
		// Get certain plugin info
		cursor := collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": pluginName})
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

func (s *Storage) CreateOnboarding(request CreateSettingsRequest, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")

	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})
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

	// Add new plugin to onboarding here
	settings := SettingsDB{
		ID:             userId,
		EnabledPlugins: request.EnabledPlugins,
		Meditation:     request.Meditation,
		Finance:        request.Finance,
		Elevator:       request.Elevator,
	}

	result, err := collection.InsertOne(ctx, settings)

	if err != nil {
		return result.InsertedID.(string), err
	}
	return "Inserted", err
}
func (s *Storage) CreatePluginSettings(request interface{}, userId string, pluginName string, ctx context.Context) (string, error) {
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
				EnabledPlugins: append(settingsRecord.EnabledPlugins, PluginName(pluginName)),
			}

			switch pluginName {
			case "finance":
				newSettings.Finance = *(request.(*FinanceSettings))
				newSettings.Meditation = settingsRecord.Meditation
				newSettings.Elevator = settingsRecord.Elevator
			case "meditation":
				newSettings.Meditation = *(request.(*MeditationSettings))
				newSettings.Finance = settingsRecord.Finance
				newSettings.Elevator = settingsRecord.Elevator
			case "elevator":
				newSettings.Elevator = *(request.(*ElevatorSettings))
				newSettings.Finance = settingsRecord.Finance
				newSettings.Meditation = settingsRecord.Meditation
			}

			// Update the settings with the new settings
			result := collection.FindOneAndUpdate(ctx, bson.M{"_id": userId}, bson.M{"$set": newSettings})
			if result.Err() != nil {
				return "", result.Err()
			} else {
				return "Created", nil
			}
		}

	} else {
		// Create new settings since the user has no settings yet
		if validateSettings(request) != nil {
			return "Invalid settings", errors.New("Invalid settings: " + validateSettings(request).Error())
		}
		settings := SettingsDB{
			ID:             userId,
			EnabledPlugins: []PluginName{PluginName(pluginName)},
		}

		switch pluginName {
		case "finance":
			settings.Finance = *(request.(*FinanceSettings))
		case "meditation":
			settings.Meditation = *(request.(*MeditationSettings))
		case "elevator":
			settings.Elevator = *(request.(*ElevatorSettings))
		}

		// Insert the settings
		_, err := collection.InsertOne(ctx, settings)
		if err != nil {
			return "", err
		}
		return "Inserted", err
	}
}

func (s *Storage) UpdatePluginSettings(request interface{}, userId string, pluginName string, ctx context.Context) (string, error) {
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
			newSettings.Elevator = settingsRecord.Elevator
		case "meditation":
			newSettings.Meditation = *(request.(*MeditationSettings))
			newSettings.Finance = settingsRecord.Finance
			newSettings.Elevator = settingsRecord.Elevator
		case "elevator":
			newSettings.Elevator = *(request.(*ElevatorSettings))
			newSettings.Finance = settingsRecord.Finance
			newSettings.Meditation = settingsRecord.Meditation

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

func (s *Storage) Delete(userId string, ctx context.Context, plugin string) (interface{}, error) {
	collection := s.db.Collection("settings")

	// Check if user exists
	userResult := collection.FindOne(ctx, bson.M{"_id": userId})
	if userResult.Err() != nil {
		return nil, errors.New("User not found")
	}

	// Delete the plugin settings for the specified user
	if plugin == "" {
		_, err := collection.DeleteOne(ctx, bson.M{"_id": userId})
		if err != nil {
			return nil, err
		}
	} else {
		// validate plugin name
		if !isValidPlugins(plugin) {
			return errors.New("Error"), errors.New("Invalid plugin name")
		}
		pluginName := PluginName(plugin)
		// Delete the plugin from the enabledPlugins array and its associated saved values
		update := bson.M{
			"$pull":  bson.M{"enabledPlugins": pluginName},
			"$unset": bson.M{plugin: ""},
		}

		_, err := collection.UpdateOne(ctx, bson.M{"_id": userId}, update)
		if err != nil {
			return nil, err
		}

		// Query the user document after the update
		var updatedUserSettings SettingsDB
		err = userResult.Decode(&updatedUserSettings)
		if err != nil {
			return nil, err
		}

		// If no plugins are left, delete the entire user settings
		// remove plugin from enabledPlugins array
		if len(updatedUserSettings.EnabledPlugins)-1 == 0 {
			_, err := collection.DeleteOne(ctx, bson.M{"_id": userId})
			if err != nil {
				return nil, err
			}
		}
	}

	return "Deleted", nil
}

// function to validate each setting of a plugin
// Todo: make dynamic
func validateSettings(request interface{}) error {
	newSettings := SettingsDB{}

	switch request.(type) {
	case *FinanceSettings:
		newSettings.Finance = *(request.(*FinanceSettings))
		if !isValidNotificationType(newSettings.Finance.PeriodNotifications) || !isValidStrategy(newSettings.Finance.Strategy) {
			return errors.New("Invalid notification type, either strategy or notification")
		}

	case *MeditationSettings:
		newSettings.Meditation = *(request.(*MeditationSettings))
		if !isValidNotificationType(newSettings.Meditation.PeriodNotifications) {
			return errors.New("Invalid notification type")
		}

	case *CreateSettingsRequest:
		req := request.(*CreateSettingsRequest)
		newSettings.Finance = req.Finance
		newSettings.Meditation = req.Meditation
		newSettings.EnabledPlugins = req.EnabledPlugins
		newSettings.Elevator = req.Elevator
		if !isValidNotificationType(newSettings.Finance.PeriodNotifications) || !isValidStrategy(newSettings.Finance.Strategy) || !isValidNotificationType(newSettings.Meditation.PeriodNotifications) || !isValidPlugins(newSettings.EnabledPlugins) {
			return errors.New("Invalid settings, notification, strategy, or plugin not supported")
		}

	case FinanceSettings:
		newSettings.Finance = request.(FinanceSettings)
		if !isValidNotificationType(newSettings.Finance.PeriodNotifications) || !isValidStrategy(newSettings.Finance.Strategy) {
			return errors.New("Invalid notification type, either strategy or notification")
		}

	case MeditationSettings:
		newSettings.Meditation = request.(MeditationSettings)
		if !isValidNotificationType(newSettings.Meditation.PeriodNotifications) {
			return errors.New("Invalid notification type")
		}

	case CreateSettingsRequest:
		req := request.(CreateSettingsRequest)
		newSettings.Finance = req.Finance
		newSettings.Meditation = req.Meditation
		newSettings.EnabledPlugins = req.EnabledPlugins
		newSettings.Elevator = req.Elevator
		if !isValidNotificationType(newSettings.Finance.PeriodNotifications) || !isValidStrategy(newSettings.Finance.Strategy) || !isValidNotificationType(newSettings.Meditation.PeriodNotifications) || !isValidPlugins(newSettings.EnabledPlugins) {
			return errors.New("Invalid settings, notification, strategy, or plugin not supported")
		}
	case ElevatorSettings:
		newSettings.Elevator = request.(ElevatorSettings)
		if !isValidNotificationType(newSettings.Elevator.PeriodNotifications) {
			return errors.New("Invalid notification type")
		}
	case *ElevatorSettings:
		newSettings.Elevator = *(request.(*ElevatorSettings))
		if !isValidNotificationType(newSettings.Elevator.PeriodNotifications) {
			return errors.New("Invalid notification type")
		}

	default:
		return errors.New("Invalid request type")
	}

	if !isValidPlugins(newSettings.EnabledPlugins) {
		return errors.New("Invalid plugin")
	}

	return nil
}

func isValidPlugins(plugin interface{}) bool {
	switch reflect.TypeOf(plugin).Kind() {
	case reflect.String:
		plugin := PluginName(plugin.(string))
		return validPlugins[plugin]
	case reflect.Slice:
		plugins := plugin.([]PluginName)
		for _, p := range plugins {
			if !validPlugins[p] {
				return false
			}
		}
		return true
	default:
		return false
	}
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
