package model

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/duyk16/notifying-server/storage"
)

type Notification struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    string             `json:"userId" bson:"userId"`
	Message   string             `json:"message" bson:"message"`
	IsRead    bool               `json:"isRead" bson:"isRead"`
	IsSMS     bool               `json:"isSMS" bson:"isSMS"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func SaveAllNotifications(users []string, message string) (err error) {
	if len(users) == 0 {
		log.Printf("All users are zero, can't save any message")
		return nil
	}

	var notifications []interface{}
	for _, userId := range users {
		notify := Notification{
			ID:        primitive.NewObjectID(),
			UserId:    userId,
			Message:   message,
			IsRead:    false,
			IsSMS:     false,
			CreatedAt: time.Now(),
		}
		notifications = append(notifications, notify)
	}
	_, err = storage.Notifications.InsertMany(context.Background(), notifications)
	return err
}

func SaveNotification(userId, message string) (err error) {
	notify := Notification{
		ID:        primitive.NewObjectID(),
		UserId:    userId,
		Message:   message,
		IsRead:    false,
		IsSMS:     false,
		CreatedAt: time.Now(),
	}
	_, err = storage.Notifications.InsertOne(context.Background(), notify)
	return err
}

func ReadNotification(userId string, id primitive.ObjectID) (err error) {
	res := storage.Notifications.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{"isRead": true},
		},
		options.FindOneAndUpdate().SetUpsert(false),
	)

	return res.Err()
}

func GetNotifications(userId string) (err error, notifications []Notification, unReadCount int64) {
	unReadCount, err = storage.Notifications.CountDocuments(
		context.Background(),
		bson.M{"userId": userId},
	)

	if err != nil {
		return err, notifications, unReadCount
	}

	cur, err := storage.Notifications.Find(
		context.Background(),
		bson.M{"userId": userId},
		options.Find().SetProjection(bson.M{
			"_id":       1,
			"message":   1,
			"isRead":    1,
			"createdAt": 1,
		}),
	)

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var notify Notification
		err = cur.Decode(&notify)
		if err != nil {
			return err, notifications, unReadCount
		}
		notifications = append(notifications, notify)
	}

	return nil, notifications, unReadCount
}
