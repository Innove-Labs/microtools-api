package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // For loading .env configuration
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file -->")
	}

	// MongoClient = initMongoDB()
	// initRedis()
	//go updateApiHitCounts()

	log.Println("Database initialized")
	router := mux.NewRouter()

	// will implement the count increment worker later
	//router.Use(APICounterMiddleware)

	// router.HandleFunc("/api/v1/user/register", RegisterUserHandler).Methods("POST")

	// auth required apis
	// router.Handle("/api/v1/email/validate", JWTAuthMiddleware(http.HandlerFunc(ValidateEmailHandler))).Methods("POST")
	router.Handle("/api/v1/email/validate", http.HandlerFunc(ValidateEmailHandler)).Methods("POST")


	// Start server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
