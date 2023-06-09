package progress

import (
	"cmd/http/main.go/internal/settings"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
)

type NotificationType string

type Experience map[settings.PluginName]float64
type ExperienceToNewLevel map[settings.PluginName]float64

const maxLevel = 6

const experienceToNewLevel = 50

type Db struct {
	// A list with the Plugins that the user has enabled.
	ID         string     `json:"id" bson:"_id"`
	Experience Experience `json:"experience" bson:"experience"`
}

type Response struct {
	Experience           Experience           `json:"level"`
	ExperienceToNewLevel ExperienceToNewLevel `json:"experienceToNewLevel"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Get(userId string, ctx context.Context) (Response, error) {
	collection := s.db.Collection("progress")
	userCollection := s.db.Collection("users")

	// Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := userResult.Err(); err != nil {
		return Response{}, errors.New("User not found!")
	}

	var db Db
	var err error

	err = collection.FindOne(ctx, bson.M{"_id": userId}).Decode(&db)
	if err != nil {
		return Response{}, err
	}

	// Calculate level
	level := make(Experience)
	toNewLevel := make(ExperienceToNewLevel)

	for plugin, experience := range db.Experience {
		calculatedLevel := math.Floor(experience / float64(experienceToNewLevel))
		if calculatedLevel > maxLevel {
			calculatedLevel = maxLevel
		}
		level[plugin] = float64(calculatedLevel)
		toNewLevel[plugin] = math.Mod(experience, float64(experienceToNewLevel))
	}
	return Response{Experience: level,
		ExperienceToNewLevel: toNewLevel}, nil
}

func (s *Storage) AddExperience(userId string, ctx context.Context, plugin settings.PluginName, experienceToAdd float64) error {
	collection := s.db.Collection("progress")
	userCollection := s.db.Collection("users")

	// Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})
	if err := userResult.Err(); err != nil {
		// Create user if not exists
		_, err := userCollection.InsertOne(ctx, bson.M{"_id": userId})
		if err != nil {
			return err
		}
	}

	var db Db
	err := collection.FindOne(ctx, bson.M{"_id": userId}).Decode(&db)
	if err != nil {
		// Create user settings if not existsy
		db = Db{
			ID:         userId,
			Experience: make(Experience),
		}
		_, err = collection.InsertOne(ctx, db)
		if err != nil {
			return err
		}
	}
	// Add experience to the plugin

	db.Experience[plugin] += experienceToAdd

	// Update the user settings in the database
	_, err = collection.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": db})
	if err != nil {
		return err
	}

	return nil
}
