package router

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/innovelabs/microtools-go/internal/handlers"
	"github.com/innovelabs/microtools-go/internal/middleware"
	"github.com/innovelabs/microtools-go/internal/services/generator"
)

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

// SetupRouter configures and returns the application router
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.APICounterMiddleware)

	// API routes
	router.Handle("/api/v1/validate/email", http.HandlerFunc(handlers.ValidateEmailHandler)).Methods("POST")
	router.Handle("/api/v1/validate/ip", http.HandlerFunc(handlers.ValidateIPHandler)).Methods("POST")
	router.Handle("/api/v1/validate/iban", http.HandlerFunc(handlers.ValidateIBANHandler)).Methods("POST")
	router.Handle("/api/v1/generate/qr", http.HandlerFunc(handlers.QRHandler)).Methods("POST")
	barcodeSvc := generator.NewDefaultBarcodeService()
	router.Handle("/api/v1/generate/barcode", handlers.GenerateBarcodeHandler(barcodeSvc)).Methods("POST")

	// Public APIs
	router.Handle("/api/v1/live", http.HandlerFunc(handlers.LiveHandler)).Methods("GET")

	// Parse templates
	homeTmpl := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/pages/home.html"))
	emailTmpl := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/pages/email.html"))
	ipTmpl := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/pages/ip.html"))
	ibanTmpl := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/pages/iban.html"))
	qrTmpl := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/pages/qr.html"))
	barcodeTmpl := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/pages/barcode.html"))

	// UI routes
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

	return router
}
