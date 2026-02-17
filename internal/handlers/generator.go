package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/innovelabs/microtools-go/internal/models"
	"github.com/innovelabs/microtools-go/internal/services/generator"
)

// QRHandler handles QR code generation requests
func QRHandler(w http.ResponseWriter, r *http.Request) {
	var req models.QRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	png, err := generator.GenerateQR(req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(png)
}

// GenerateBarcodeHandler handles barcode generation requests
func GenerateBarcodeHandler(barcodeSvc generator.BarcodeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.GenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		data, contentType, err := barcodeSvc.Generate(req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
