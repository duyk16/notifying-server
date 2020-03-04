package server

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/duyk16/notifying-server/config"
	u "github.com/duyk16/notifying-server/util"

	"github.com/gorilla/websocket"
)

var connections = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)               // broadcast channel
var upgrader = websocket.Upgrader{}
var usageCores = 1

func Init() {
	// Set max CPU cores
	var threads = config.ServerConfig.Threads
	if threads < 0 {
		log.Printf("Threads must be greater or equal 0")
		return
	} else if threads == 0 {
		threads = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(threads)
	log.Printf("Server running on %v CPU cores", threads)

	// Allow CORs
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// Hanler routers
	http.HandleFunc("/ws", socketHandler)
	http.HandleFunc("/message", messageHandler)

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	connections[ws] = false
	log.Printf("Connect from IP: %v \n Total connections: %v", ws.RemoteAddr().String(), len(connections))

	// Handle timeout on Authentication
	go func() {
		time.Sleep(time.Duration(config.ServerConfig.AuthTimeout) * time.Second)
		if valid, ok := connections[ws]; ok {
			if !valid {
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
		case "auth":
			if msg.Token != "" {
				connections[ws] = true
				ws.WriteJSON(u.T{
					"status": "ok",
				})
			} else {
				ws.WriteJSON(u.T{
					"status": "error",
					"error":  "Your token is not valid",
				})
				return
			}
		}

		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		u.JSON(w, 400, u.T{
			"status": "error",
			"error":  "Body error",
		})
		return
	}

	switch request.Type {
	case "NORMAL":
		{
			u.JSON(w, 400, u.T{
				"status": "ok",
			})
		}
	default:
		{

			u.JSON(w, 400, u.T{
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
