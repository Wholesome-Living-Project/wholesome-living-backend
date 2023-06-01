package meditation

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"time"
)

type MeditationDB struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	UserID         string             `json:"userId" bson:"userId"`
	MeditationTime int                `json:"meditationTime" bson:"meditationTime"`
	EndTime        int64              `json:"endTime" bson:"endTime"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Create(request createMeditationRequest, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("meditation")

	createdAt := time.Now().Unix()

	meditation := MeditationDB{
		ID:             primitive.NewObjectID(),
		UserID:         userId,
		MeditationTime: request.MeditationTime,
		EndTime:        createdAt,
	}

	result, err := collection.InsertOne(ctx, meditation)

	if err != nil {
		return "", err
	}

	// convert the object id to a string
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Storage) Get(meditationID string, ctx context.Context) (MeditationDB, error) {
	collection := s.db.Collection("meditation")
	meditationRecord := MeditationDB{}

	objectID, err := primitive.ObjectIDFromHex(meditationID)
	if err != nil {
		return meditationRecord, err
	}

	cursor := collection.FindOne(ctx, bson.M{"_id": objectID})

	if err := cursor.Err(); err != nil {
		return meditationRecord, err
	}

	if err := cursor.Decode(&meditationRecord); err != nil {
		return meditationRecord, err
	}
	return meditationRecord, nil
}

func (s *Storage) GetAllOfOneUserBetweenTimeAndDuration(userId string, times map[string]int64, ctx context.Context) ([]MeditationDB, error) {
	// get all meditations of one user between two times
	collection := s.db.Collection("meditation")
	var cursor *mongo.Cursor
	var err error
	if times["endTime"] == 0 {
		times["endTime"] = time.Now().Unix()
	}
	if times["durationEnd"] == 0 {
		times["durationEnd"] = math.MaxInt64
	}
	meditations := make([]MeditationDB, 0)
	cursor, err = collection.Find(ctx, bson.M{"userId": userId, "endTime": bson.M{"$gte": times["startTime"], "$lte": times["endTime"]}, "meditationTime": bson.M{"$gte": times["startDuration"], "$lte": times["durationEnd"]}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var meditation MeditationDB
		if err := cursor.Decode(&meditation); err != nil {
			return nil, err
		}
		meditations = append(meditations, meditation)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// return the meditation list
	return meditations, nil
}
