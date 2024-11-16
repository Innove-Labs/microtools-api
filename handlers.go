package main

import (
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type EmailRequest struct {
	Email string `json:"email"`
}

type UserRequest struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Company string `json:"company"`
	Country string `json:"country"`
}

func ValidateEmailHandler(w http.ResponseWriter, r *http.Request) {
	var email EmailRequest

	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Validating email: ", email.Email)
	emailValidationResult := ValidateEmail(email.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"validationResult": emailValidationResult})
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user UserRequest

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
	jwt, jwtErr := GenerateJWT(user.Email)
	if jwtErr != nil {
		http.Error(w, jwtErr.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully", "token": jwt})
}
