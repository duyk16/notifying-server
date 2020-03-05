package server

import (
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"

	"github.com/duyk16/notifying-server/config"
)

var AllUsers []string // all users of service

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
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Hanler routers
	r := mux.NewRouter()
	r.HandleFunc("/ws", SocketHandler)
	r.HandleFunc("/messages", MessageHandler).Methods("POST")

	// Handler subcribe on redis
	SubscribeRedisChannel(config.ServerConfig.Database.RedisChannel)

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":"+config.ServerConfig.Port, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
