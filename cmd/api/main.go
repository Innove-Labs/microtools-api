package main

import (
	"log"
	"net/http"

	"github.com/innovelabs/microtools-go/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize databases (currently commented out)
	// cfg := config.LoadConfig()
	// handlers.MongoClient = database.InitMongoDB()
	// redisClient := database.InitRedis()

	log.Println("Database initialized")

	// Setup router
	r := router.SetupRouter()

	// Start server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
