package Settings_Test

import (
	"context"
	"testing"

	"cmd/http/main.go/internal/settings"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestStorage_Get(t *testing.T) {
	// Set up a test MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	assert.NoError(t, err)
	defer client.Disconnect(context.Background())

	// Prepare test data
	db := client.Database("testdb")
	settingsCollection := db.Collection("settings")
	usersCollection := db.Collection("users")

	// Insert a user document
	userID := "user1"
	_, err = usersCollection.InsertOne(context.Background(), bson.M{"_id": userID})
	assert.NoError(t, err)

	// Insert a settings document
	settingsData := bson.M{
		"_id":            userID,
		"enabledPlugins": []string{"finance", "meditation"},
		"meditation": bson.M{
			"meditationTime":      10,
			"notifications":       true,
			"amountNotifications": 5,
			"periodNotifications": "Day",
		},
		"finance": bson.M{
			"notifications":       true,
			"amountNotifications": 3,
			"periodNotifications": "Month",
			"strategy":            "Round",
			"strategyAmount":      100,
			"investmentGoal":      5000,
			"investmentTimeGoal":  365,
		},
	}
	_, err = settingsCollection.InsertOne(context.Background(), settingsData)
	assert.NoError(t, err)

	// Create a new storage instance for testing
	storage := settings.NewStorage(db)

	// Test case 1: Get all settings
	result, err := storage.Get(userID, context.Background(), "")
	assert.NoError(t, err)
	assert.Equal(t, "user1", result.ID)
	assert.ElementsMatch(t, []settings.PluginName{"finance", "meditation"}, result.EnabledPlugins)
	assert.Equal(t, 10, result.Meditation.MeditationTimeGoal)
	assert.Equal(t, true, result.Meditation.Notifications)
	assert.Equal(t, 5, result.Meditation.AmountNotifications)
	assert.Equal(t, settings.NotificationTypeDay, result.Meditation.PeriodNotifications)
	assert.Equal(t, true, result.Finance.Notifications)
	assert.Equal(t, 3, result.Finance.AmountNotifications)
	assert.Equal(t, settings.NotificationTypeMonth, result.Finance.PeriodNotifications)
	assert.Equal(t, settings.StrategyTypeRound, result.Finance.Strategy)
	assert.Equal(t, 100, result.Finance.StrategyAmount)
	assert.Equal(t, 5000, result.Finance.InvestmentGoal)
	assert.Equal(t, 365, result.Finance.InvestmentTimeGoal)

	// Test case 2: Get finance plugin settings
	result, err = storage.Get(userID, context.Background(), "finance")
	assert.NoError(t, err)
	assert.Equal(t, "user1", result.ID)
	assert.ElementsMatch(t, []settings.PluginName{"finance", "meditation"}, result.EnabledPlugins)
	assert.Equal(t, true, result.Finance.Notifications)
	assert.Equal(t, 3, result.Finance.AmountNotifications)
	assert.Equal(t, settings.NotificationTypeMonth, result.Finance.PeriodNotifications)
	assert.Equal(t, settings.StrategyTypeRound, result.Finance.Strategy)
	assert.Equal(t, 100, result.Finance.StrategyAmount)
	assert.Equal(t, 5000, result.Finance.InvestmentGoal)
	assert.Equal(t, 365, result.Finance.InvestmentTimeGoal)

	// Test case 3: Get settings for invalid user
	_, err = storage.Get("invalidUser", context.Background(), "")
	assert.Error(t, err)
	assert.EqualError(t, err, "User not found!")

	// Test case 4: Get settings for invalid plugin
	_, err = storage.Get(userID, context.Background(), "invalidPlugin")
	assert.Error(t, err)
	assert.EqualError(t, err, "Plugin not found!")
}

// TODO mocking of the mongo db
//
//func TestCreateOnboarding(t *testing.T) {
//	db := getMockDB()
//	storage := settings.NewStorage(db)
//
//	request := settings.CreateSettingsRequest{
//		EnabledPlugins: []settings.PluginName{"finance", "meditation"},
//		Meditation: settings.MeditationSettings{
//			MeditationTimeGoal:  20,
//			Notifications:       true,
//			AmountNotifications: 2,
//			PeriodNotifications: settings.NotificationTypeWeek,
//		},
//		Finance: settings.FinanceSettings{
//			Notifications:       true,
//			AmountNotifications: 3,
//			PeriodNotifications: settings.NotificationTypeMonth,
//			Strategy:            settings.StrategyTypeRound,
//			StrategyAmount:      100,
//			InvestmentGoal:      5000,
//			InvestmentTimeGoal:  365,
//		},
//	}
//
//	// Test case 1: Create onboarding settings for a new user
//	result, err := storage.CreateOnboarding(request, "user1", context.Background())
//	assert.NoError(t, err)
//	assert.Equal(t, "Inserted", result)
//
//	// Test case 2: Create onboarding settings for an existing user
//	result, err = storage.CreateOnboarding(request, "user2", context.Background())
//	assert.NoError(t, err)
//	assert.Equal(t, "", result)
//	assert.EqualError(t, err, "User already has onboarding settings")
//
//	// Test case 3: Create onboarding settings for a non-existing user
//	result, err = storage.CreateOnboarding(request, "user3", context.Background())
//	assert.EqualError(t, err, "User not found!")
//}
//
//type MockCollection struct {
//	mock.Mock
//}
//
//func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
//	args := m.Called(ctx, filter, opts)
//	return args.Get(0).(*mongo.SingleResult)
//}
//
//func (m *MockCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
//	args := m.Called(ctx, filter, update, opts)
//	return args.Get(0).(*mongo.SingleResult)
//}
//
//func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	args := m.Called(ctx, filter, update, opts)
//	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
//}
//
//func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
//	args := m.Called(ctx, document, opts)
//	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
//}
//
//func TestCreatePluginSettings(t *testing.T) {
//	db := getMockDB()
//	storage := settings.NewStorage(db)
//
//	request := settings.MeditationSettings{
//		MeditationTimeGoal:  20,
//		Notifications:       true,
//		AmountNotifications: 2,
//		PeriodNotifications: settings.NotificationTypeWeek,
//	}
//
//	// Test case 1: Create meditation plugin settings for a new user
//	result, err := storage.CreatePluginSettings(&request, "user1", "meditation", context.Background())
//	assert.NoError(t, err)
//	assert.Equal(t, "Inserted", result)
//
//	// Test case 2: Create finance plugin settings for an existing user
//	financeRequest := settings.FinanceSettings{
//		Notifications:       true,
//		AmountNotifications: 3,
//		PeriodNotifications: settings.NotificationTypeMonth,
//		Strategy:            settings.StrategyTypeRound,
//		StrategyAmount:      100,
//		InvestmentGoal:      5000,
//		InvestmentTimeGoal:  365,
//	}
//	result, err = storage.CreatePluginSettings(&financeRequest, "user2", "finance", context.Background())
//	assert.NoError(t, err)
//	assert.Equal(t, "Inserted", result)
//
//	// Test case 3: Create finance plugin settings for a user with existing settings
//	result, err = storage.CreatePluginSettings(&financeRequest, "user2", "finance", context.Background())
//	assert.EqualError(t, err, "User already has finance settings, use put")
//}
//
//func TestUpdatePluginSettings(t *testing.T) {
//	db := getMockDB()
//	storage := settings.NewStorage(db)
//
//	request := settings.FinanceSettings{
//		Notifications:       true,
//		AmountNotifications: 5,
//		PeriodNotifications: settings.NotificationTypeWeek,
//		Strategy:            settings.StrategyTypePercent,
//		StrategyAmount:      10,
//		InvestmentGoal:      10000,
//		InvestmentTimeGoal:  730,
//	}
//
//	// Test case 1: Update finance plugin settings for an existing user
//	result, err := storage.UpdatePluginSettings(&request, "user1", "finance", context.Background())
//	assert.NoError(t, err)
//	assert.Equal(t, "Updated", result)
//
//	// Test case 2: Update finance plugin settings for a non-existing user
//	result, err = storage.UpdatePluginSettings(&request, "user3", "finance", context.Background())
//	assert.EqualError(t, err, "User not found!")
//
//	// Test case 3: Update non-existing plugin settings for an existing user
//	result, err = storage.UpdatePluginSettings(&request, "user1", "meditation", context.Background())
//	assert.EqualError(t, err, "User does not have meditation settings, use create")
//
//	// Test case 4: Update plugin settings with an invalid plugin name
//	result, err = storage.UpdatePluginSettings(&request, "user1", "invalid", context.Background())
//	assert.EqualError(t, err, "Invalid plugin name")
//}
//
//func TestGetPluginSettings(t *testing.T) {
//	db := getMockDB()
//	storage := settings.NewStorage(db)
//
//	// Test case 1: Get meditation plugin settings for an existing user
//	result, err := storage.Get("user1", context.Background(), "meditation")
//	assert.NoError(t, err)
//	assert.Equal(t, &settings.MeditationSettings{
//		MeditationTimeGoal:  20,
//		Notifications:       true,
//		AmountNotifications: 2,
//		PeriodNotifications: settings.NotificationTypeWeek,
//	}, result)
//
//	// Test case 2: Get finance plugin settings for a non-existing user
//	result, err = storage.Get("user3", context.Background(), "finance")
//	assert.EqualError(t, err, "User not found!")
//
//	// Test case 3: Get non-existing plugin settings for an existing user
//	result, err = storage.Get("user1", context.Background(), "fitness")
//	assert.EqualError(t, err, "User does not have fitness settings")
//}
//
//func TestDeletePluginSettings(t *testing.T) {
//	db := getMockDB()
//	storage := settings.NewStorage(db)
//
//	// Test case 1: Delete finance plugin settings for an existing user
//	result, err := storage.Delete("user1", context.Background(), "finance")
//	assert.NoError(t, err)
//	assert.Equal(t, "Deleted", result)
//
//	// Test case 2: Delete plugin settings for a non-existing user
//	result, err = storage.Delete("user3", context.Background(), "finance")
//	assert.EqualError(t, err, "User not found!")
//
//	// Test case 3: Delete non-existing plugin settings for an existing user
//	result, err = storage.Delete("user1", context.Background(), "meditation")
//	assert.EqualError(t, err, "User does not have meditation settings")
//}
