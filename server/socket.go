package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/duyk16/notifying-server/config"
	"github.com/duyk16/notifying-server/model"
	"github.com/duyk16/notifying-server/util"
)

var AllUsers []string // all users of service

var upgrader = websocket.Upgrader{}
var connections = make(map[*websocket.Conn]string)
var clients = make(map[string]*websocket.Conn)

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

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
	handleMessage(ws)
}

func handleMessage(ws *websocket.Conn) {
	// Make sure we close the connection when the function return
	defer delete(connections, ws)
	defer ws.Close()

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Socket close: %v \n Total connections: %v", err.Error(), len(connections))
			return
		}

		switch msg.Type {
		case "AUTH":
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

			clients[msg.Token] = ws
		case "GET_DATA":
			userId := connections[ws]
			if userId == "" {
				ws.WriteJSON(util.T{
					"status": "error",
					"error":  "Unauthentication",
				})
				return
			}

			err, notifications, unReadCount := model.GetNotifications(userId)

			log.Printf("User get data %v, %v", notifications, unReadCount)

			if err != nil {
				ws.WriteJSON(util.T{
					"status": "error",
					"error":  "Try again",
				})
			}

			ws.WriteJSON(util.T{
				"status": "ok",
				"data": bson.M{
					"notifications": notifications,
					"unReadCount":   unReadCount,
				},
			})
		case "READ":
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
			} else {
				ws.WriteJSON(util.T{
					"status": "ok",
				})
			}

			model.ReadNotification(userId, id)
		default:
			ws.WriteJSON(util.T{
				"status": "error",
				"error":  "Type is not valid",
			})
			return
		}
	}
}

func SentMessage(ID string, message string) {
	conn, ok := clients[ID]
	if ok {
		conn.WriteJSON(message)
	}

	model.SaveNotification(ID, message)
}

func SentMessageToAll(message string) {
	for _, conn := range clients {
		conn.WriteJSON(message)
	}

	model.SaveAllNotifications(AllUsers, message)
}
