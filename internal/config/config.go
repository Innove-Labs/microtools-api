package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds all the configuration variables.
type Config struct {
	MongoURI      string
	RedisURI      string
	JWTSecret     string
	CounterApiKey string
}

// LoadConfig loads the environment variables from .env file and returns a Config object.
func LoadConfig() *Config {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Retrieve the variables from the environment
	return &Config{
		MongoURI:      os.Getenv("MONGO_URI"),
		RedisURI:      os.Getenv("REDIS_URI"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		CounterApiKey: os.Getenv("COUNTER_API_KEY"),
	}
}
