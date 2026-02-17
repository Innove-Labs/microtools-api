package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/innovelabs/microtools-go/internal/models"
	"github.com/innovelabs/microtools-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

// RegisterUserHandler handles user registration requests
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	collection := MongoClient.Database("microapps").Collection("users")
	existingUser, err := collection.FindOne(r.Context(), bson.M{"email": user.Email}).Raw()
	if existingUser != nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}
	_, err = collection.InsertOne(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jwt, jwtErr := utils.GenerateJWT(user.Email)
	if jwtErr != nil {
		http.Error(w, jwtErr.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully", "token": jwt})
}
