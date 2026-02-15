package main

import (
	"html/template"
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

	router.Use(APICounterMiddleware)

	// router.HandleFunc("/api/v1/user/register", RegisterUserHandler).Methods("POST")

	// router.Handle("/api/v1/email/validate", JWTAuthMiddleware(http.HandlerFunc(ValidateEmailHandler))).Methods("POST")
	router.Handle("/api/v1/validate/email", http.HandlerFunc(ValidateEmailHandler)).Methods("POST")
	router.Handle("/api/v1/validate/ip", http.HandlerFunc(ValidateIPHandler)).Methods("POST")
	router.Handle("/api/v1/generate/qr", http.HandlerFunc(QRHandler)).Methods("POST")

	// public apis
	router.Handle("/api/v1/live", http.HandlerFunc(LiveHandler)).Methods("GET")

	tmpl, err := template.ParseGlob("views/layout.html")
	if err != nil {
		log.Fatalf("Error parsing layout template: %v", err)
	}
	tmpl, err = tmpl.ParseGlob("views/partials/*.html")
	if err != nil {
		log.Fatalf("Error parsing partial templates: %v", err)
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, nil); err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}).Methods("GET")

	// Start server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
