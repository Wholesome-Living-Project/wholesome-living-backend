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

// AND ADD ALSO HERE
// map of PluginName to Capitized names
var pluginNameToCap = map[PluginName]string{
	PluginNameFinance:    "Finance",
	PluginNameMeditation: "Meditation",
	PluginNameElevator:   "Elevator",
}

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
	getName() string
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
	NotificationId      int              `json:"notificationId" bson:"notificationId"`
}

func (f FinanceSettings) getPeriodNotifications() NotificationType {
	return f.PeriodNotifications
}

func (f FinanceSettings) getName() string {
	return "finance"
}

// TODO check if enough
func (f FinanceSettings) validate() error {
	if !isValidNotificationType(f) || !isValidStrategy(f.Strategy) {
		return errors.New("invalid finance strategy")
	}
	return nil
}

type MeditationSettings struct {
	MeditationTimeGoal  int              `json:"meditationTimeGoal" bson:"meditationTimeGoal"`
	Notifications       bool             `json:"notifications" bson:"notifications"`
	AmountNotifications int              `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications NotificationType `json:"periodNotifications" bson:"periodNotifications"`
	NotificationId      int              `json:"notificationId" bson:"notificationId"`
}

func (m MeditationSettings) getPeriodNotifications() NotificationType {
	return m.PeriodNotifications
}

func (m MeditationSettings) getName() string {
	return "meditation"
}

// TODO check if enough
func (m MeditationSettings) validate() error {
	if !isValidNotificationType(m) {
		return errors.New("Invalid notification type")
	}
	return nil
}

type ElevatorSettings struct {
	Notifications       bool             `json:"notifications" bson:"notifications"`
	AmountNotifications int              `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications NotificationType `json:"periodNotifications" bson:"periodNotifications"`
	Goal                int              `json:"goal" bson:"goal"`
	NotificationId      int              `json:"notificationId" bson:"notificationId"`
}

func (e ElevatorSettings) getPeriodNotifications() NotificationType {
	return e.PeriodNotifications
}

func (e ElevatorSettings) getName() string {
	return "elevator"
}

// TODO check if enough
func (e ElevatorSettings) validate() error {
	if !isValidNotificationType(e) {
		return errors.New("Invalid notification type")
	}
	return nil
}

// SettingsDB is the struct that is stored in the database
// enabledPlugis -> array of plugins of user
// TypeSetting -> user settings
type SettingsDB struct {
	ID             string             `json:"id" bson:"_id"`
	EnabledPlugins []PluginName       `json:"enabledPlugins" bson:"enabledPlugins"`
	Meditation     MeditationSettings `json:"meditation" bson:"meditation,omitempty" `
	Finance        FinanceSettings    `json:"finance" bson:"finance,omitempty"`
	Elevator       ElevatorSettings   `json:"elevator" bson:"elevator,omitempty"`
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
func (s *Storage) CreateOnboarding(request CreateSettingsRequest, userId string, ctx context.Context) (string, error) {
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
	if err := validateSettingsRequest(request); err != nil {
		return "Invalid settings", err
	}

	// Create settings
	settings, err := createEnabledSettings(request, userId)
	if err != nil {
		return "", err
	}

	// Insert settings
	result, err := collection.InsertOne(ctx, settings)
	if err != nil {
		return result.InsertedID.(string), err
	}

	return "Created", err
}

func (s *Storage) CreatePluginSettings(request SingleSetting, userId string, ctx context.Context) error {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")
	pluginName := request.getName()

	// Check if user exists
	user := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := user.Err(); err != nil {
		return errors.New("User not found!")
	}

	// Check if plugin exists
	if err := validateSettings(request); err != nil {
		return err
	}

	// Check if user already has onboarding settings
	settings := collection.FindOne(ctx, bson.M{"_id": userId})
	if settings.Err() != nil {
		// Create new settings since the user has no settings yet
		sett := SettingsDB{
			ID:             userId,
			EnabledPlugins: []PluginName{PluginName(pluginName)},
		}

		// Insert the settings
		if _, err := collection.InsertOne(ctx, sett); err != nil {
			return err
		}
	}

	// settings should exist now
	settings = collection.FindOne(ctx, bson.M{"_id": userId})

	// Create plugin settings and keep the other settings
	var settingsRecord SettingsDB
	if err := settings.Decode(&settingsRecord); err != nil {
		return err
	}

	// Check if user already has the specified plugin settings
	for _, plugin := range settingsRecord.EnabledPlugins {
		if plugin == PluginName(pluginName) {
			return errors.New("User already has " + pluginName + " settings")
		}
	}

	// Add new plugin to onboarding here
	updatedEnabled := append(settingsRecord.EnabledPlugins, PluginName(pluginName))

	// Update the settings with the new settings
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": userId},
		bson.M{"$set": bson.M{
			"enabledPlugins": updatedEnabled,
			pluginName:       request,
		},
		})

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (s *Storage) UpdatePluginSettings(request SingleSetting, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")
	pluginName := request.getName()

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
	if !isValidPlugins(plugin) {
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

func validateSettingsRequest(request CreateSettingsRequest) error {
	if !isValidPlugins(request.EnabledPlugins) {
		return errors.New("Invalid enabled plugin ")
	}

	for _, v := range request.EnabledPlugins {
		pluginCap, ok := pluginNameToCap[PluginName(v)]
		if !ok {
			return errors.New("Invalid plugin name: " + string(v))
		}

		// Validate field if present
		field := reflect.ValueOf(request).FieldByName(string(pluginCap))
		if field.CanInterface() {
			log.Println("Validating field: ", field.Type())
			singset, ok := field.Interface().(SingleSetting)
			if !ok {
				return errors.New("Invalid field type: " + string(pluginCap) + " is not a SingleSetting")
			}

			if err := singset.validate(); err != nil {
				return err
			}
		}

	}
	return nil
}

func createEnabledSettings(request CreateSettingsRequest, userId string) (SettingsDB, error) {
	settingsDB := SettingsDB{ID: userId, EnabledPlugins: request.EnabledPlugins}

	for _, v := range request.EnabledPlugins {
		pluginCap, ok := pluginNameToCap[PluginName(v)]
		if !ok {
			return settingsDB, errors.New("Invalid plugin name: " + string(v))
		}

		if srcField := reflect.ValueOf(request).FieldByName(string(pluginCap)); srcField.CanInterface() {
			log.Println("Validating src-field: ", srcField.Type())
			dstField := reflect.ValueOf(&settingsDB).Elem().FieldByName(string(pluginCap))
			dstField.Set(srcField)
		}
	}

	return settingsDB, nil
}

// function to validate each setting of a plugin
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
