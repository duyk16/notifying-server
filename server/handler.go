package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/duyk16/notifying-server/model"
	"github.com/duyk16/notifying-server/storage"
	"github.com/duyk16/notifying-server/util"
)

var broadcast = make(chan Message) // broadcast channel

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": "error",
			"error":  "Body error",
		})
		return
	}

	err = handleRequestMessage(request)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": "error",
			"error":  "Input data is not valid",
		})
		return
	}

	util.JSON(w, 200, util.T{"status": "ok"})
}

func SubscribeRedisChannel(channel string) {
	pubsub, err := storage.SubscribeChanel(channel)

	if err != nil {
		log.Printf("Fail to sub channel %v", channel)
		os.Exit(1)
		return
	}

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var request Request
			err := json.Unmarshal([]byte(msg.Payload), &request)

			if err != nil {
				log.Printf("Fail to decode message\n Channel: %v\n Message: %v", msg.Channel, msg.Payload)
			}

			err = handleRequestMessage(request)

			if err != nil {
				log.Printf("Fail to handle message %v", request)
			}
		}
	}()

}

func handleRequestMessage(request Request) error {
	switch request.Type {
	case "NORMAL":
		{
			for _, userId := range request.Users {
				SentMessage(userId, request.Message)
				model.SaveNotification(userId, request.Message)

			}
			return nil
		}
	case "ALL":
		{
			SentMessageToAll(request.Message)
			model.SaveAllNotifications(AllUsers, request.Message)
			return nil
		}
	default:
		{

			return errors.New("Input data is not valid")
		}
	}
}
