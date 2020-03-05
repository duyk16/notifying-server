package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Threads     int      `json:"threads"`
	Name        string   `json:"name"`
	Port        string   `json:"port"`
	Database    Database `json:"database"`
	AuthTimeout int      `json:"auth_timeout"`
}

type Database struct {
	MongoURL      string `json:"mongo_url"`
	MongoDB       string `json:"mongo_db"`
	RedisURL      string `json:"redis_url"`
	RedisPoolSize int    `json:"redis_pool_size"`
	RedisPassword string `json:"redis_password"`
	RedisDB       int    `json:"redis_db"`
	RedisChannel  string `json:"redis_channel"`
}

var ServerConfig Config

func Init() {
	configFileName := "config.json"
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("File error: ", err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&ServerConfig); err != nil {
		log.Fatal("Config error: ", err.Error())
	}
}
