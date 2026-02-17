package database

import (
	"context"
	"log"
	"time"

	"github.com/innovelabs/microtools-go/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDB initializes MongoDB client
func InitMongoDB() *mongo.Client {
	cfg := config.LoadConfig()
	log.Println("Initializing MongoDB... with uri: ", cfg.MongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return client
}
