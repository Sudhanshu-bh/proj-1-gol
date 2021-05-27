package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Email struct {
	User string
}

type Details struct {
	User     string `bson:"user"`
	Password string `bson:"password"`
}

type ChangePass struct {
	User     string `bson:"user"`
	CurrPass string `bson:"password"`
	NewPass  string `bson:"newpass"`
}

var DBResult Details

// var uri string = "mongodb://localhost:27017"

var uri string = "mongodb+srv://admin1:hpcadmin55@cluster0.cgpy4.mongodb.net/users?retryWrites=true&w=majority"

func CheckUserInDB(email string) (error, error) {
	// In the function signature, first error is for the error which DB query will return
	// and second error is the database connectivity error.

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("users").Collection("users")

	err = collection.FindOne(context.TODO(), bson.M{"username": email}).Decode(&DBResult)

	err2 := client.Disconnect(ctx)

	if err2 != nil {
		return nil, err2
	}
	fmt.Println("Connection to MongoDB closed.")

	return err, nil
}

func ChangePassInDB(creds ChangePass) (mongo.UpdateResult, error, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return mongo.UpdateResult{}, nil, err
	}

	// Check the connection
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return mongo.UpdateResult{}, nil, err
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("users").Collection("users")

	filter := bson.M{"username": creds.User, "password": creds.CurrPass}

	UpdateOneStruct, err := collection.UpdateOne(
		context.TODO(),
		filter,
		bson.D{
			{"$set", bson.M{"password": creds.NewPass}},
		})

	err2 := client.Disconnect(ctx)

	if err2 != nil {
		return *UpdateOneStruct, nil, err2
	}
	fmt.Println("Connection to MongoDB closed.")

	return *UpdateOneStruct, err, nil
}
