package storage

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/duyk16/notifying-server/config"
)

var Database *mongo.Database
var Notifications *mongo.Collection

func ConnectMongo() {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.ServerConfig.Database.MongoURL))
	if err != nil {
		log.Fatal(err)
	}

	// Connect the mongo client to the MongoDB server
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)

	// Ping MongoDB
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Printf("Could not ping to mongo db service: %v\n", err)
		os.Exit(1)
		return
	}

	log.Println("Connected to MongoDB")

	Database = client.Database(config.ServerConfig.Database.MongoDB)
	Notifications = Database.Collection("notifications")
	createIndexes()
}

func createIndexes() {
	Notifications.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"userId": 1},
		Options: options.Index().SetUnique(true),
	})
}
