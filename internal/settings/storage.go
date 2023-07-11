package settings

import (
	"context"
	"errors"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

type SingleSetting interface {
	getPeriodNotifications() NotificationType
	validate() error
}

// Interest rate of investment TODO maybe there can be multiple different interest rates for different investmens
type FinanceSettings struct {
	Notifications       bool             `json:"notifications" bson:"notifications"`
	AmountNotifications int              `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications NotificationType `json:"periodNotifications" bson:"periodNotifications"`
	Strategy            StrategyType     `json:"strategy" bson:"strategy"`
	StrategyAmount      int              `json:"strategyAmount" bson:"strategyAmount"`
	InvestmentGoal      int              `json:"investmentGoal" bson:"investmentGoal"`
	InvestmentTimeGoal  int              `json:"investmentTimeGoal" bson:"investmentTimeGoal"`
}

func (f FinanceSettings) getPeriodNotifications() NotificationType {
	return f.PeriodNotifications
}

func (f FinanceSettings) validate() error {
	if !isValidNotificationType(f) || !isValidStrategy(f.Strategy) {
		return errors.New("Invalid notification type, either strategy or notification")
	}
	return nil
}

type MeditationSettings struct {
	MeditationTimeGoal  int              `json:"meditationTimeGoal" bson:"meditationTimeGoal"`
	Notifications       bool             `json:"notifications" bson:"notifications"`
	AmountNotifications int              `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications NotificationType `json:"periodNotifications" bson:"periodNotifications"`
}

func (m MeditationSettings) getPeriodNotifications() NotificationType {
	return m.PeriodNotifications
}

// TODO
func (f MeditationSettings) validate() error {
	return nil
}

type ElevatorSettings struct {
	Notifications       bool             `json:"notifications" bson:"notifications"`
	AmountNotifications int              `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications NotificationType `json:"periodNotifications" bson:"periodNotifications"`
	Goal                int              `json:"goal" bson:"goal"`
}

func (e ElevatorSettings) getPeriodNotifications() NotificationType {
	return e.PeriodNotifications
}

// TODO
func (f ElevatorSettings) validate() error {
	return nil
}

// SettingsDB is the struct that is stored in the database
// enabledPlugis -> array of plugins of user
// TypeSetting -> user settings
type SettingsDB struct {
	ID             string             `json:"id" bson:"_id"`
	EnabledPlugins []PluginName       `json:"enabledPlugins" bson:"enabledPlugins"`
	Meditation     MeditationSettings `json:"meditation" bson:"meditation"`
	Finance        FinanceSettings    `json:"finance" bson:"finance"`
	Elevator       ElevatorSettings   `json:"elevator" bson:"elevator"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Get(userId string, plugin string, ctx context.Context) (SettingsDB, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")
	settingsRecord := SettingsDB{}

	// Check if user exists
	user := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := user.Err(); err != nil {
		return settingsRecord, errors.New("User not found!")
	}

	// No plugin - Get all plugins
	if plugin == "" {
		cursor := collection.FindOne(ctx, bson.M{"_id": userId})
		if err := cursor.Err(); err != nil {
			return settingsRecord, err
		}

		// Decode the record
		if err := cursor.Decode(&settingsRecord); err != nil {
			return settingsRecord, err
		}

		return settingsRecord, nil
	}

	// Check if plugin exists
	pluginName := PluginName(plugin)
	if _, ok := validPlugins[pluginName]; !ok {
		return settingsRecord, errors.New("Plugin not found!")
	}

	// Get certain plugin info
	cursor := collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": pluginName})
	if err := cursor.Err(); err != nil {
		return settingsRecord, err
	}

	// Decode the record
	if err := cursor.Decode(&settingsRecord); err != nil {
		return settingsRecord, err
	}

	return settingsRecord, nil
}

// TODO: Settings should be when onboarding is made or a user choses a ned plugin
func (s *Storage) CreateOnboarding(request CreateOnboardingSettingResponse, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")

	// Check if user exists
	user := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := user.Err(); err != nil {
		return "User not found", err
	}

	// Check if user already has onboarding settings
	userSettings := collection.FindOne(ctx, bson.M{"_id": userId})
	if userSettings.Err() == nil {
		return "", errors.New("User already has onboarding settings")
	}

	// Validate request
	// FIXME: loop over all enabled plugins and check if they are valid
	if err := validateSettingsRequest(request); err != nil {
		return "Invalid settings", err
	}

	// FIXME: make more dynamic
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

	return "", err
}

func (s *Storage) CreatePluginSettings(request SingleSetting, userId string, pluginName string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")

	// Check if user exists
	user := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := user.Err(); err != nil {
		return "", errors.New("User not found!")
	}

	// Check if user already has onboarding settings
	settings := collection.FindOne(ctx, bson.M{"_id": userId})
	if settings.Err() == nil {

		// Check if user already has the specified plugin settings
		if collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": pluginName}).Err() == nil {
			return "", errors.New("User already has " + pluginName + " settings, use put")
		}

		// Create plugin settings and keep the other settings
		var settingsRecord SettingsDB
		if err := settings.Decode(&settingsRecord); err != nil {
			return "", err
		}

		// Check if plugin exists
		if err := validateSettings(request); err != nil {
			return "", err
		}

		// Add new plugin to onboarding here
		settingsRecord.EnabledPlugins = append(settingsRecord.EnabledPlugins, PluginName(pluginName))

		// FIXME: make more dynamic
		// maybe through update logic or throught setters on settings
		switch pluginName {
		case "finance":
			settingsRecord.Finance = *(request.(*FinanceSettings))
		case "meditation":
			settingsRecord.Meditation = *(request.(*MeditationSettings))
		case "elevator":
			settingsRecord.Elevator = *(request.(*ElevatorSettings))
		}

		// Update the settings with the new settings
		result := collection.FindOneAndReplace(ctx, bson.M{"_id": userId}, bson.M{"$set": settingsRecord})
		if result.Err() != nil {
			return "", result.Err()
		}

		return "Created", nil

	} else {
		// Create new settings since the user has no settings yet
		if err := validateSettings(request); err != nil {
			return "Invalid settings", errors.New("Invalid settings: " + err.Error())
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

// TODO: encapsulate the single setting in update request object
func (s *Storage) UpdatePluginSettings(request SingleSetting, userId string, pluginName string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")

	// Check if user already has the specified plugin settings
	oldSettings := collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": pluginName})
	if oldSettings.Err() != nil {
		return "", errors.New("User does not have " + pluginName + " settings")
	}

	// Update plugin settings
	var settingsRecord SettingsDB
	if err := oldSettings.Decode(&settingsRecord); err != nil {
		return "", err
	}

	// Validate the request updates
	if err := validateSettings(request); err != nil {
		return "", err
	}

	// Update the OnboardingSettings with the new settings
	cursor := collection.FindOneAndUpdate(ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{pluginName: request}})
	if cursor.Err() != nil {
		return "", cursor.Err()
	}

	// TODO: return the updated document
	// Decode the record
	//  if err := cursor.Decode(&settingsRecord); err != nil {
	//		return settingsRecord, err
	// }

	// Fixme: see above
	return "Updated", nil
}

func (s *Storage) Delete(userId string, plugin string, ctx context.Context) error {
	collection := s.db.Collection("settings")

	// Check if user exists
	userResult := collection.FindOne(ctx, bson.M{"_id": userId})
	if userResult.Err() != nil {
		log.Println(userResult.Err())
		return errors.New("No plugin-settings found for user")
	}

	// Delete all settings, if no plugin is specified
	if plugin == "" {
		_, err := collection.DeleteOne(ctx, bson.M{"_id": userId})
		if err != nil {
			return err
		}
		return nil
	}

	// Validate plugin name
	if isValidPlugins(plugin) != true {
		return errors.New("Invalid plugin name")
	}

	// Load settings
	var sdb SettingsDB
	err := userResult.Decode(&sdb)
	if err != nil {
		return err
	}

	// Check if plugin is enabled
	enabled := false
	for _, v := range sdb.EnabledPlugins {
		if v == PluginName(plugin) {
			enabled = true
			break
		}
	}

	// nothing to do if plugin is not enabled, success
	if !enabled {
		return nil
	}

	// Delete the plugin from the enabledPlugins array and its associated saved values
	pluginName := PluginName(plugin)
	update := bson.M{
		"$pull":  bson.M{"enabledPlugins": pluginName},
		"$unset": bson.M{plugin: ""},
	}

	// Update the settings
	_, err = collection.UpdateOne(ctx, bson.M{"_id": userId}, update)
	if err != nil {
		return err
	}

	// If no plugins are left, delete the entire user settings
	// remove plugin from enabledPlugins array
	if len(sdb.EnabledPlugins)-1 == 0 {
		_, err := collection.DeleteOne(ctx, bson.M{"_id": userId})
		if err != nil {
			return err
		}
	}

	return nil
}

func validateSettingsRequest(request CreateOnboardingSettingResponse) error {
	if !isValidNotificationType(request.Finance) || !isValidStrategy(request.Finance.Strategy) || !isValidNotificationType(request.Meditation) || !isValidPlugins(request.EnabledPlugins) {
		return errors.New("Invalid settings, notification, strategy, or plugin not supported")
	}

	return nil
}

// function to validate each setting of a plugin
// Todo: make dynamic
func validateSettings(setting SingleSetting) error {
	if err := setting.validate(); err != nil {
		return err
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

func isValidNotificationType(setting SingleSetting) bool {
	notifications := setting.getPeriodNotifications()
	switch notifications {
	case NotificationTypeDay, NotificationTypeWeek, NotificationTypeMonth:
		return true
	default:
		return false
	}
}

func isValidStrategy(strat StrategyType) bool {
	switch strat {
	case StrategyTypePercent, StrategyTypeRound, StrategyTypePlus:
		return true
	default:
		return false
	}
}
