package server

import (
	"github.com/duyk16/notifying-server/model"
	"github.com/gorilla/websocket"
)

var AllUsers []string // all users of service

var clients map[string]*websocket.Conn

func CreateClient(ID string, ws *websocket.Conn) {
	clients[ID] = ws
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
