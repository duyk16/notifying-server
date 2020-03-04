package server

import (
	"encoding/json"
	"log"
	"net/http"

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

	switch request.Type {
	case "NORMAL":
		{
			for _, userId := range request.Users {
				SentMessage(userId, request.Message)
			}
			util.JSON(w, 200, util.T{"status": "ok"})
			return
		}
	case "ALL":
		{
			SentMessageToAll(request.Message)
			util.JSON(w, 200, util.T{"status": "ok"})
			return
		}
	default:
		{

			util.JSON(w, 400, util.T{
				"status": "error",
				"error":  "Input data is not valid",
			})
			return
		}
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		for client := range connections {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(connections, client)
			}
		}
	}
}
