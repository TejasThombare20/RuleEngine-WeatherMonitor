package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSetup() *mongo.Client {

	// MONGODB_URI := os.Getenv("MONGODB_URI")
	MONGODB_URI := "mongodb+srv://tejas:tejas@cluster0.o2ucgqs.mongodb.net/zeotap?retryWrites=true&w=majority"

	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	fmt.Println("Connected to MongoDB server")

	return client

}

var Client *mongo.Client = DBSetup()

var RuleCollection = Client.Database("zeotap").Collection("rules")

// func getCollection(client *mongo.Collection, collectionName string) *mongo.Collection {

// 	var collection *mongo.Collection = client.Database("zeotap").Collection(collectionName)

// 	return collection

// }
