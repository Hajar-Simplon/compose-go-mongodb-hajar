package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	connectionString = fmt.Sprintf("mongodb://%s:%s",
		os.Getenv("MONGODB_HOST"),
		os.Getenv("MONGODB_PORT"),
	)
)

func main() {
	ctx := context.Background()

	client, _ := mongo.Connect(options.Client().ApplyURI(connectionString))

	res, err := client.
		Database("test").
		Collection("test").
		InsertOne(ctx, bson.M{"created_at": time.Now().Format(time.RFC3339)})
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
