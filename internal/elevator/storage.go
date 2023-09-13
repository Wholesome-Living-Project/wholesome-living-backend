package elevator

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ElevatorDB struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       string             `json:"userId" bson:"userId"`
	Time         int64              `json:"time" bson:"time"`
	Stairs       bool               `json:"stairs" bson:"stairs"`
	AmountStairs int                `json:"amountStairs" bson:"amountStairs"`
	HeightGain   int64              `json:"heightGain" bson:"heightGain"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Create(request CreateElevatorRequest, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("elevator")

	createdAt := time.Now().Unix()
	if request.AmountStairs != 0 && !request.Stairs {
		return "", fmt.Errorf("amountStairs can only be set if stairs is true")
	}

	elevator := ElevatorDB{
		ID:           primitive.NewObjectID(),
		UserID:       userId,
		Time:         createdAt,
		Stairs:       request.Stairs,
		AmountStairs: request.AmountStairs,
		HeightGain:   request.HeightGain,
	}

	result, err := collection.InsertOne(ctx, elevator)

	if err != nil {
		return "", err
	}

	// convert the object id to a string
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Storage) Get(elevatorID string, ctx context.Context) (ElevatorDB, error) {
	collection := s.db.Collection("elevator")
	elevatorRecord := ElevatorDB{}

	objectID, err := primitive.ObjectIDFromHex(elevatorID)
	if err != nil {
		return elevatorRecord, err
	}

	cursor := collection.FindOne(ctx, bson.M{"_id": objectID})

	if err := cursor.Err(); err != nil {
		return elevatorRecord, err
	}

	if err := cursor.Decode(&elevatorRecord); err != nil {
		return elevatorRecord, err
	}
	return elevatorRecord, nil
}

func (s *Storage) GetAllOfOneUserBetweenTimeAndDuration(userId string, times map[string]int64, gain map[string]int64, ctx context.Context) ([]ElevatorDB, error) {
	// get all elevators of one user between two times
	collection := s.db.Collection("elevator")
	var cursor *mongo.Cursor
	var err error
	if times["endTime"] == 0 {
		times["endTime"] = time.Now().Unix()
	}
	if times["durationEnd"] == 0 {
		times["durationEnd"] = math.MaxInt64
	}
	if gain["maxGain"] == 0 {
		gain["maxGain"] = math.MaxInt64
	}

	elevators := make([]ElevatorDB, 0)
	cursor, err = collection.Find(ctx, bson.M{"userId": userId, "time": bson.M{"$gte": times["startTime"], "$lte": times["endTime"]}, "amountStairs": bson.M{"$gte": times["durationStart"], "$lte": times["durationEnd"]}, "heightGain": bson.M{"$gte": gain["minGain"], "$lte": gain["maxGain"]}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var elevator ElevatorDB
		if err := cursor.Decode(&elevator); err != nil {
			return nil, err
		}
		elevators = append(elevators, elevator)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// return the elevator list
	return elevators, nil
}
