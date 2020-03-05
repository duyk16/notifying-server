package storage

import (
	"log"
	"os"

	"github.com/go-redis/redis"

	"github.com/duyk16/notifying-server/config"
)

var client *redis.Client

func ConnectRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     config.ServerConfig.Database.RedisURL,
		PoolSize: config.ServerConfig.Database.RedisPoolSize,
		Password: config.ServerConfig.Database.RedisPassword,
		DB:       config.ServerConfig.Database.RedisDB,
	})
	_, err := client.Ping().Result()

	if err != nil {
		log.Printf("Fail to connect Redis at %v", config.ServerConfig.Database.RedisURL)
		os.Exit(1)
		return
	}

	log.Println("Connected to Redis")
}

func SubscribeChanel(channel string) (*redis.PubSub, error) {
	pubsub := client.Subscribe(channel)
	_, err := pubsub.Receive()
	if err != nil {
		log.Printf("Subsribe channel %v fail", channel)
		os.Exit(1)
	}

	return pubsub, err
}

func PublistChannel(channel, message string) error {
	// Publish a message.
	err := client.Publish(channel, message).Err()
	return err
}
