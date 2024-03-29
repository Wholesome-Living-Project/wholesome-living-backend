package finance

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type financeDB struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       string             `json:"userId" bson:"userId"`
	SpendingTime int64              `json:"spendingTime" bson:"spendingTime"`
	Amount       float64            `json:"amount" bson:"amount"`
	Saving       float64            `json:"saving" bson:"saving"`
	Description  string             `json:"description" bson:"description"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) create(request CreateSpendingRequest, userId string, ctx context.Context) (string, error) {
	collection := s.db.Collection("investment")
	userCollection := s.db.Collection("users")

	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userId})

	if err := userResult.Err(); err != nil {
		return "", err
	}
	statement := financeDB{
		ID:           primitive.NewObjectID(),
		UserID:       userId,
		SpendingTime: request.SpendingTime,
		Amount:       request.Amount,
		Saving:       request.Saving,
		Description:  request.Description,
	}

	result, err := collection.InsertOne(ctx, statement)

	if err != nil {
		return "", err
	}

	// convert the object id to a string
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Storage) get(investmentID string, ctx context.Context) (financeDB, error) {
	collection := s.db.Collection("investment")
	db := financeDB{}

	objectID, err := primitive.ObjectIDFromHex(investmentID)
	if err != nil {
		return db, err
	}

	cursor := collection.FindOne(ctx, bson.M{"_id": objectID})

	if err := cursor.Err(); err != nil {
		return db, err
	}

	if err := cursor.Decode(&db); err != nil {
		return db, err
	}
	return db, nil
}

func (s *Storage) getAllOfOneUser(userID string, ctx context.Context) ([]financeDB, error) {
	collection := s.db.Collection("investment")
	userCollection := s.db.Collection("users")

	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userID})

	if err := userResult.Err(); err != nil {
		fmt.Println("Error finding user:", err)
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	investments := make([]financeDB, 0)
	for cursor.Next(ctx) {
		var investment financeDB
		if err := cursor.Decode(&investment); err != nil {
			return nil, err
		}
		investments = append(investments, investment)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// return the investment list
	return investments, nil
}

func (s *Storage) getAllOfOneUserBetweenTime(id string, startTime int64, endTime int64, ctx context.Context) ([]financeDB, error) {
	// get all investments of one user between two times
	collection := s.db.Collection("investment")
	var cursor *mongo.Cursor
	var err error
	// different query if endtime is 0
	if endTime == 0 {
		cursor, err = collection.Find(ctx, bson.M{"userId": id, "investmentTime": bson.M{"$gte": startTime}})
	} else {
		cursor, err = collection.Find(ctx, bson.M{"userId": id, "investmentTime": bson.M{"$gte": startTime, "$lte": endTime}})
	}
	if err != nil {
		return nil, err
	}

	investments := make([]financeDB, 0)
	for cursor.Next(ctx) {
		var investment financeDB
		if err := cursor.Decode(&investment); err != nil {
			return nil, err
		}
		investments = append(investments, investment)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// return the investment list
	return investments, nil
}
