package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/innovelabs/microtools-go/internal/models"
	"github.com/innovelabs/microtools-go/internal/services/validation"
)

// ValidateEmailHandler handles email validation requests
func ValidateEmailHandler(w http.ResponseWriter, r *http.Request) {
	var email models.EmailRequest

	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Validating email: ", email.Email)
	formattedEmail := strings.TrimSpace(email.Email)
	emailValidationResult := validation.ValidateEmail(formattedEmail)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"validationResult": emailValidationResult})
}

// ValidateIPHandler handles IP validation/geolocation requests
func ValidateIPHandler(w http.ResponseWriter, r *http.Request) {
	var ip models.IPRequest

	err := json.NewDecoder(r.Body).Decode(&ip)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	log.Println("Validating IP: ", ip.IP)
	formattedIP := strings.TrimSpace(ip.IP)
	ipValidationResult, err := validation.ValidateIP(formattedIP)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"validationResult": ipValidationResult})
}

// ValidateIBANHandler handles IBAN validation requests
func ValidateIBANHandler(w http.ResponseWriter, r *http.Request) {
	var ibanReq models.IBANRequest

	err := json.NewDecoder(r.Body).Decode(&ibanReq)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	log.Println("Validating IBAN:", ibanReq.IBAN)
	formattedIBAN := strings.TrimSpace(ibanReq.IBAN)
	ibanValidationResult := validation.ValidateIBAN(formattedIBAN)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"validationResult": ibanValidationResult,
	})
}
