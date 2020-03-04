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
	Database    Database `json:"database"`
	AuthTimeout int      `json:"auth_timeout"`
}

type Database struct {
	MongoURL      string `json:"mongo_url"`
	RedisURL      string `json:"redis_url"`
	RedisUser     string `json:"redis_user"`
	RedisPassword string `json:"redis_password"`
	RedisDB       int    `json:"redis_db"`
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
