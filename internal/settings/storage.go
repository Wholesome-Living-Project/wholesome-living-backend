package settings

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeditationSettings struct {
	// The user's meditation time goal.
	MeditationTimeGoal  int                    `json:"meditationTime" bson:"meditationTime"`
	Notifications       bool                   `json:"notifications" bson:"notifications"`
	AmountNotifications int                    `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications enumNotificationPeriod `json:"periodNotifications" bson:"periodNotifications"`
}

type enumNotificationPeriod struct {
	Day   bool `json:"day" bson:"day"`
	Week  bool `json:"week" bson:"week"`
	Month bool `json:"month" bson:"month"`
}
type enumStrategy struct {
	Round   bool `json:"round" bson:"round"`
	Plus    bool `json:"plus" bson:"plus"`
	Percent bool `json:"percent" bson:"percent"`
}

type FinanceSettings struct {
	Notifications       bool                   `json:"notifications" bson:"notifications"`
	AmountNotifications int                    `json:"amountNotifications" bson:"amountNotifications"`
	PeriodNotifications enumNotificationPeriod `json:"periodNotifications" bson:"periodNotifications"`
	Strategy            enumStrategy           `json:"strategy" bson:"strategy"`
	StrategyAmount      int                    `json:"strategyAmount" bson:"strategyAmount"`
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
		// TODO maybe check if plugin exists
		// Get certain plugin info
		cursor := collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": plugin})
		if err := cursor.Err(); err != nil {
			return settingsRecord, err
		}
		// TODO return only relevant info from selected plugin
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
	// TODO maybe overwrite the settings if the user already has them --> check
	if err := userResult.Err(); err != nil {
		return "User not found", err
	}
	settings := SettingsDB{
		ID:             userId,
		EnabledPlugins: request.EnabledPlugins,
		Meditation: MeditationSettings{
			MeditationTimeGoal:  request.Meditation.MeditationTimeGoal,
			Notifications:       request.Meditation.Notifications,
			AmountNotifications: request.Meditation.AmountNotifications,
			PeriodNotifications: enumNotificationPeriod{
				Day:   request.Meditation.PeriodNotifications.Day,
				Week:  request.Meditation.PeriodNotifications.Week,
				Month: request.Meditation.PeriodNotifications.Month,
			},
		},
		Finance: FinanceSettings{
			Notifications:       request.Finance.Notifications,
			AmountNotifications: request.Finance.AmountNotifications,
			PeriodNotifications: enumNotificationPeriod{
				Day:   request.Finance.PeriodNotifications.Day,
				Week:  request.Finance.PeriodNotifications.Week,
				Month: request.Finance.PeriodNotifications.Month,
			},
			Strategy: enumStrategy{
				Round:   request.Finance.Strategy.Round,
				Plus:    request.Finance.Strategy.Plus,
				Percent: request.Finance.Strategy.Percent,
			},
			StrategyAmount:     request.Finance.StrategyAmount,
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

func (s *Storage) createFinanceSettings(request FinanceSettings, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")
	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := userResult.Err(); err != nil {
		return "", errors.New("User not found!")
	}
	// Check if user already has onboarding settings
	settings := collection.FindOne(ctx, bson.M{"_id": userId})
	if settings.Err() == nil {
		// check if user already has finance settings
		if collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": "finance"}).Err() == nil {
			return "", errors.New("User already has finance settings, use put")
		} else {
			// create finance settings and keep the other settings
			var settingsRecord SettingsDB
			if err := settings.Decode(&settingsRecord); err != nil {
				return "", err
			}
			newSettings := SettingsDB{
				ID:             userId,
				EnabledPlugins: append(settingsRecord.EnabledPlugins, "finance"),
				Finance:        request,
				// take rest from old settings
				Meditation: settingsRecord.Meditation,
			}
			// update the settings with the new settings
			result := collection.FindOneAndUpdate(ctx, bson.M{"_id": userId}, bson.M{"$set": newSettings})
			if result.Err() != nil {
				return "", result.Err()
			} else {
				return "Updated", nil
			}
		}

	} else {
		// create new settings, since user has no settings yet
		settings := SettingsDB{
			ID:             userId,
			EnabledPlugins: []string{"finance"},
			Finance:        request,
		}
		// insert the settings
		_, err := collection.InsertOne(ctx, settings)
		if err != nil {
			return "", err
		}
		return "Inserted", err
	}

	return "Somethin went wrong", nil
}
func (s *Storage) createMeditationSettings(request MeditationSettings, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("settings")
	userCollection := s.db.Collection("users")
	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := userResult.Err(); err != nil {
		return "", errors.New("User not found!")
	}
	// Check if user already has onboarding settings
	settings := collection.FindOne(ctx, bson.M{"_id": userId})
	if settings.Err() == nil {
		// check if user already has finance settings
		if collection.FindOne(ctx, bson.M{"_id": userId, "enabledPlugins": "meditation"}).Err() == nil {
			return "", errors.New("User already has meditation settings, use put")
		} else {
			// create finance settings and keep the other settings
			var settingsRecord SettingsDB
			if err := settings.Decode(&settingsRecord); err != nil {
				return "", err
			}
			newSettings := SettingsDB{
				ID: userId,
				// add meditation to the enabled plugins
				EnabledPlugins: append(settingsRecord.EnabledPlugins, "meditation"),
				Meditation:     request,
				// take rest from old settings
				Finance: settingsRecord.Finance,
			}
			// update the settings with the new settings
			result := collection.FindOneAndUpdate(ctx, bson.M{"_id": userId}, bson.M{"$set": newSettings})
			if result.Err() != nil {
				return "", result.Err()
			} else {
				return "Updated", nil
			}
		}

	} else {
		// create new settings, since user has no settings yet
		settings := SettingsDB{
			ID:             userId,
			EnabledPlugins: []string{"finance"},
			Meditation:     request,
		}
		// insert the settings
		_, err := collection.InsertOne(ctx, settings)
		if err != nil {
			return "", err
		}
		return "Inserted", err
	}

	return "Somethin went wrong", nil
}
