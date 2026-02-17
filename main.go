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

type PageData struct {
	Title       string
	Description string
	Canonical   string
}

func renderPage(tmpl *template.Template, data PageData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

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
	router.Handle("/api/v1/validate/iban", http.HandlerFunc(ValidateIBANHandler)).Methods("POST")
	router.Handle("/api/v1/generate/qr", http.HandlerFunc(QRHandler)).Methods("POST")

	barcodeSvc := NewDefaultBarcodeService()
	router.Handle("/api/v1/generate/barcode", http.HandlerFunc(GenerateBarcodeHandler(barcodeSvc))).Methods("POST")

	// public apis
	router.Handle("/api/v1/live", http.HandlerFunc(LiveHandler)).Methods("GET")

	// Page templates
	homeTmpl := template.Must(template.ParseFiles("views/base.html", "views/pages/home.html"))
	emailTmpl := template.Must(template.ParseFiles("views/base.html", "views/pages/email.html"))
	ipTmpl := template.Must(template.ParseFiles("views/base.html", "views/pages/ip.html"))
	ibanTmpl := template.Must(template.ParseFiles("views/base.html", "views/pages/iban.html"))
	qrTmpl := template.Must(template.ParseFiles("views/base.html", "views/pages/qr.html"))
	barcodeTmpl := template.Must(template.ParseFiles("views/base.html", "views/pages/barcode.html"))

	router.HandleFunc("/", renderPage(homeTmpl, PageData{
		Title:       "Micro API - Free Developer APIs for Email, IP, QR & Barcode",
		Description: "Free REST APIs for email validation, IP geolocation, QR code generation, and barcode generation. Simple JSON interface, no API key required.",
		Canonical:   "/",
	})).Methods("GET")

	router.HandleFunc("/email-validation-api", renderPage(emailTmpl, PageData{
		Title:       "Free Email Validation API - Syntax, Domain & Disposable Check",
		Description: "Validate email addresses with syntax checking, domain verification, MX record lookup, and disposable email detection. Free REST API with JSON response.",
		Canonical:   "/email-validation-api",
	})).Methods("GET")

	router.HandleFunc("/ip-geolocation-api", renderPage(ipTmpl, PageData{
		Title:       "Free IP Geolocation API - Country, City & Timezone Lookup",
		Description: "Look up any IP address to get country, region, city, coordinates, and timezone. Free REST API powered by MaxMind GeoIP2.",
		Canonical:   "/ip-geolocation-api",
	})).Methods("GET")

	router.HandleFunc("/iban-validation-api", renderPage(ibanTmpl, PageData{
		Title:       "Free IBAN Validation API - Format, Checksum & Country Verification",
		Description: "Validate International Bank Account Numbers (IBAN) with comprehensive checks including format validation, mod-97 checksum verification, and country-specific rules for 60+ countries.",
		Canonical:   "/iban-validation-api",
	})).Methods("GET")

	router.HandleFunc("/qr-code-generator-api", renderPage(qrTmpl, PageData{
		Title:       "Free QR Code Generator API - Text, URL, WiFi, vCard & More",
		Description: "Generate QR codes as PNG images. Supports text, URLs, email, phone, WiFi, vCard, geo, events, and JSON. Free REST API.",
		Canonical:   "/qr-code-generator-api",
	})).Methods("GET")

	router.HandleFunc("/barcode-generator-api", renderPage(barcodeTmpl, PageData{
		Title:       "Free Barcode Generator API - UPC-A, EAN-13 & Code128",
		Description: "Generate 1D barcodes in PNG or SVG format. Supports UPC-A, EAN-13, and Code128 with optional human-readable text. Free REST API.",
		Canonical:   "/barcode-generator-api",
	})).Methods("GET")

	// Start server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
