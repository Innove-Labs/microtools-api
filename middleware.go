package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Endpoint-to-counter name mapping. Add new entries here when adding new routes.
var counterNames = map[string]string{
	"/api/v1/validate/email": "email-validate",
	"/api/v1/validate/ip":    "ip-validate",
	"/api/v1/generate/qr":    "qr-generate",
	"/api/v1/live":           "live",
}

const counterBaseURL = "https://api.counterapi.dev/v2/fawaz-sullias-team-2926"

// Shared HTTP client with a short timeout so counter calls don't linger.
var counterHTTPClient = &http.Client{Timeout: 5 * time.Second}

// APICounterMiddleware increments a CounterAPI.dev counter for each known endpoint.
// The counter call is fire-and-forget in a goroutine — the request is served
// immediately and is never blocked or delayed by the counter call.
func APICounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the request first — no delay.
		next.ServeHTTP(w, r)

		// Fire counter increment in background after the response is sent.
		if counterName, exists := counterNames[r.URL.Path]; exists {
			go incrementCounter(counterName)
		}
	})
}

func incrementCounter(counterName string) {
	apiKey := os.Getenv("COUNTER_API_KEY")
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
