package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Person struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type PersonRepository struct {
	mongoClient *mongo.Client
	collection  *mongo.Collection
}

func NewPersonRepository(mongoClient *mongo.Client) PersonRepository {
	return PersonRepository{
		mongoClient: mongoClient,
		collection: mongoClient.Database("mydb").Collection("people"),
	}
}

func genCtx() context.Context {
	ctx, _ :=  context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func (pr *PersonRepository) GetAll() ([]Person, error) {
	ctx := genCtx()
	var people []Person
	cursor, err := pr.collection.Find(ctx, bson.M{})
	if err != nil {
		return people, fmt.Errorf(err.Error())
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		return people, fmt.Errorf(err.Error())
	}
	return people, nil
}

func (pr *PersonRepository) CreateOne(p *Person) (interface{}, error) {
	ctx := genCtx()
	result, err := pr.collection.InsertOne(ctx, *p)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return result.InsertedID, nil
}
