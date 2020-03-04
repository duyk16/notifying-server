package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/duyk16/notifying-server/config"
	"github.com/duyk16/notifying-server/model"
	"github.com/duyk16/notifying-server/util"
)

var connections = make(map[*websocket.Conn]string) // connected clients
var upgrader = websocket.Upgrader{}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	connections[ws] = ""
	log.Printf("Connect from IP: %v \n Total connections: %v", ws.RemoteAddr().String(), len(connections))

	// Handle timeout on Authentication
	go func() {
		time.Sleep(time.Duration(config.ServerConfig.AuthTimeout) * time.Second)
		if userId, ok := connections[ws]; ok {
			if userId == "" {
				delete(connections, ws)
				ws.Close()
			}
		}
	}()

	// Handle message
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Socket close: %v \n Total connections: %v", err.Error(), len(connections))
			delete(connections, ws)
			return
		}

		switch msg.Type {
		case "AUTH":
			{

				if msg.Token == "" {
					ws.WriteJSON(util.T{
						"status": "error",
						"error":  "Your token is not valid",
					})
					return
				}

				// Assumming userId is token
				connections[ws] = msg.Token
				ws.WriteJSON(util.T{
					"status": "ok",
				})

				// Auth success and create client
				CreateClient(msg.Token, ws)
			}
		case "GET_DATA":
			{
				userId := connections[ws]
				if userId == "" {
					ws.WriteJSON(util.T{
						"status": "error",
						"error":  "Unauthentication",
					})
					return
				}
				err, data := model.GetNotifications(userId)

				if err != nil {
					ws.WriteJSON(util.T{
						"status": "error",
						"error":  "Try again",
					})
				}

				ws.WriteJSON(util.T{
					"status": "ok",
					"data":   data,
				})
			}
		case "READ":
			{
				userId := connections[ws]
				if userId == "" {
					ws.WriteJSON(util.T{
						"status": "error",
						"error":  "Unauthentication",
					})
					return
				}

				id, err := primitive.ObjectIDFromHex(msg.ID)

				if err != nil {
					ws.WriteJSON(util.T{
						"status": "error",
						"error":  "Id is not valid",
					})
				}

				model.ReadNotification(userId, id)
			}
		}
	}

}
