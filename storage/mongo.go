package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/duyk16/notifying-server/config"
)

var Database *mongo.Database
var Notifications *mongo.Collection

func Init() {
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(config.ServerConfig.Database.MongoURL)
	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		log.Println("Connect to MongoDB fail")
		return
	}
	log.Println("Connected to MongoDB")

	Database = client.Database(config.ServerConfig.Database.MongoDB)
	Notifications = Database.Collection("notifications")
	createIndexes()
}

func createIndexes() {
	Notifications.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
}
