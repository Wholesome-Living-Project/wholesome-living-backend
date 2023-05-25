package settings

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeditationSettings struct {
	// The user's meditation time goal.
	MeditationTime int `json:"meditationTime" bson:"meditationTime"`
}

type FinanceSettings struct {
	// The user's investment goal.
	InvestmentGoal int `json:"investmentGoal" bson:"investmentGoal"`
	// The user's investment time goal.
	InvestmentTimeGoal int `json:"investmentTimeGoal" bson:"investmentTimeGoal"`
	// Interest rate of investment TODO maybe there can be multiple different interest rates for different investmens
}

type SettingsDB struct {
	// A list with the Plugins that the user has enabled.
	EnabledPlugins []string `json:"enabledPlugins" bson:"enabledPlugins"`
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

func (s *Storage) Get(userId string, ctx context.Context) (SettingsDB, error) {
	collection := s.db.Collection("settings")
	settingsRecord := SettingsDB{}

	cursor := collection.FindOne(ctx, bson.M{"userId": userId})

	// TODO Get particular setting of s plugin

	if err := cursor.Err(); err != nil {
		return settingsRecord, err
	}

	if err := cursor.Decode(&settingsRecord); err != nil {
		return settingsRecord, err
	}
	return settingsRecord, nil
}

// TODO: Settings should be when onboarding is made or a user choses a ned plugin

func (s *Storage) CreateOnboarding(request createOnboardingRequest, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")

	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})

	if err := userResult.Err(); err != nil {
		return "User not found", err
	}
	settings := SettingsDB{
		ID:             userId,
		EnabledPlugins: request.EnabledPlugins,
		Meditation: MeditationSettings{
			MeditationTime: request.Meditation.MeditationTime,
		},
		Finance: FinanceSettings{
			InvestmentGoal:     request.Finance.InvestmentGoal,
			InvestmentTimeGoal: request.Finance.InvestmentTimeGoal,
		},
	}

	result, err := collection.InsertOne(ctx, settings)

	if err != nil {
		return result.InsertedID.(string), err
	}
	return "Inserted", err
}
