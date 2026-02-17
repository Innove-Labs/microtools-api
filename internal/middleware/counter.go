package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/innovelabs/microtools-go/internal/config"
)

var counterNames = map[string]string{
	"/api/v1/validate/email":   "email-validate",
	"/api/v1/validate/ip":      "ip-validate",
	"/api/v1/validate/iban":    "iban-validate",
	"/api/v1/generate/qr":      "qr-generate",
	"/api/v1/generate/barcode": "barcode-generate",
	"/api/v1/live":             "live",
}

const counterBaseURL = "https://api.counterapi.dev/v2/fawaz-sullias-team-2926"

var counterHTTPClient = &http.Client{Timeout: 5 * time.Second}

// APICounterMiddleware increments a CounterAPI.dev counter for each known endpoint
func APICounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if counterName, exists := counterNames[r.URL.Path]; exists {
			go incrementCounter(counterName)
		}
	})
}

func incrementCounter(counterName string) {
	apiKey := config.LoadConfig().CounterApiKey
	url := counterBaseURL + "/" + counterName + "/up"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[counter] failed to build request for %s: %v", counterName, err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := counterHTTPClient.Do(req)
	if err != nil {
		log.Printf("[counter] failed to call %s: %v", counterName, err)
		return
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[counter] %s returned status %d", counterName, resp.StatusCode)
		return
	}

	log.Printf("[counter] incremented %s", counterName)
}
