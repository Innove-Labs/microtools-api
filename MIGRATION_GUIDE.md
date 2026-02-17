# Go Project Restructuring - Step-by-Step Migration Guide

## Overview

This guide will help you restructure the microtools-go project from a flat structure into a proper Go Standard Project Layout with clear separation of concerns.

**Estimated Time:** 1-2 hours
**Difficulty:** Intermediate

## Prerequisites Completed ✓

- [x] Directory structure created
- [x] go.mod updated to `github.com/innovelabs/microtools-go`
- [x] Config file migrated to `internal/config/config.go`
- [x] Request models migrated to `internal/models/requests.go`
- [x] Response models migrated to `internal/models/responses.go`

---

## Phase 2: Complete Models Migration

### Step 1: Create `internal/models/iban.go`

Create the file and add the IBAN country specifications:

```bash
cat > internal/models/iban.go << 'EOF'
package models

// IBANCountrySpec defines the IBAN structure for a specific country
type IBANCountrySpec struct {
	CountryCode   string
	CountryName   string
	Length        int
	BBANFormat    string
	BankCodeStart int
	BankCodeLen   int
	AccountStart  int
	AccountLen    int
	Example       string
}

// GetIBANCountrySpecs returns the map of all supported IBAN country specifications
func GetIBANCountrySpecs() map[string]IBANCountrySpec {
	return ibanCountrySpecs
}

// ibanCountrySpecs contains IBAN specifications for 60+ countries
var ibanCountrySpecs = map[string]IBANCountrySpec{
	// Copy the entire ibanCountrySpecs map from iban-handling.go (lines 40-242)
	// SEPA Countries (European Union)
	"AD": {CountryCode: "AD", CountryName: "Andorra", Length: 24, BBANFormat: "^[0-9]{8}[A-Z0-9]{12}$",
		BankCodeStart: 4, BankCodeLen: 8, AccountStart: 12, AccountLen: 12,
		Example: "AD1200012030200359100100"},
	// ... (copy all country specs from the original file)
	"XK": {CountryCode: "XK", CountryName: "Kosovo", Length: 20, BBANFormat: "^[0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 12,
		Example: "XK051212012345678906"},
}
EOF
```

**Action:** Copy the complete `ibanCountrySpecs` map from `iban-handling.go` lines 40-242 into this file.

**Checkpoint:** Verify the file compiles:
```bash
go build ./internal/models
```

---

## Phase 3: Migrate Service Layer (Business Logic)

### Step 2: Create `internal/services/validation/email.go`

```bash
cat > internal/services/validation/email.go << 'EOF'
package validation

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/innovelabs/microtools-go/internal/models"
)

var disposableEmailDomains = []string{
	"mailinator.com",
	"10minutemail.com",
	"guerrillamail.com",
	"tempmail.net",
	"throwawaymail.com",
	"yopmail.com",
	"maildrop.cc",
	"getnada.com",
	"dispostable.com",
	"fakeinbox.com",
	"tempmail.org",
	"spamgourmet.com",
	"trashmail.com",
}

func isValidEmailSyntax(email string) bool {
	var emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func isValidDomain(domain string) bool {
	_, err := net.LookupMX(domain)
	if err == nil {
		return true
	}
	_, err = net.LookupHost(domain)
	return err == nil
}

func verifyMxRecords(email string) bool {
	domain := extractDomain(email)
	if domain == "" {
		fmt.Println("Invalid email format")
		return false
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		fmt.Println("No MX records found for domain", domain)
		return false
	}

	return true
}

func isDisposableEmail(email string) bool {
	domain := extractDomain(email)
	if domain == "" {
		fmt.Println("Invalid email format")
		return false
	}

	for _, disposableDomain := range disposableEmailDomains {
		if strings.ToLower(domain) == disposableDomain {
			return true
		}
	}
	return false
}

// ValidateEmail validates an email address with comprehensive checks
func ValidateEmail(email string) models.EmailValidation {
	emailValidationResult := models.EmailValidation{
		Email:          email,
		IsSyntaxValid:  false,
		IsDomainValid:  false,
		MxRecordsFound: false,
		IsDisposable:   false,
	}

	if isValidEmailSyntax(email) {
		emailValidationResult.IsSyntaxValid = true
	}

	domain := extractDomain(email)
	if isValidDomain(domain) {
		emailValidationResult.IsDomainValid = true
	}

	if verifyMxRecords(email) {
		emailValidationResult.MxRecordsFound = true
	}

	if isDisposableEmail(email) {
		emailValidationResult.IsDisposable = true
	}

	return emailValidationResult
}
EOF
```

**Action:** Copy the complete email validation logic from `email-handling.go`.

### Step 3: Create `internal/services/validation/ip.go`

```bash
cat > internal/services/validation/ip.go << 'EOF'
package validation

import (
	"errors"
	"log"
	"net"

	"github.com/innovelabs/microtools-go/internal/models"
	"github.com/oschwald/geoip2-golang"
)

// ValidateIP validates an IP address and returns geolocation information
func ValidateIP(ipStr string) (models.GeoIPResponse, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return models.GeoIPResponse{}, errors.New("Invalid IP address")
	}

	db, err := geoip2.Open("./assets/geolite-2-city.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	record, err := db.City(ip)
	if err != nil {
		return models.GeoIPResponse{}, errors.New("IP address not found")
	}

	resp := models.GeoIPResponse{
		IP:        ipStr,
		Country:   record.Country.Names["en"],
		Region:    "",
		City:      record.City.Names["en"],
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
		Timezone:  record.Location.TimeZone,
	}
	if len(record.Subdivisions) > 0 {
		resp.Region = record.Subdivisions[0].Names["en"]
	}

	return resp, nil
}
EOF
```

**Action:** Copy from `geolocation-handling.go`. **Note:** Update path to GeoIP database to `./assets/geolite-2-city.mmdb`.

### Step 4: Create `internal/services/validation/iban.go`

```bash
cat > internal/services/validation/iban.go << 'EOF'
package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/innovelabs/microtools-go/internal/models"
)

// Helper functions
func isIBANLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isIBANDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func formatIBANWithSpaces(iban string) string {
	formatted := ""
	for i, char := range iban {
		if i > 0 && i%4 == 0 {
			formatted += " "
		}
		formatted += string(char)
	}
	return formatted
}

func calculateMod97(numStr string) int {
	remainder := 0
	for _, digit := range numStr {
		digitVal := int(digit - '0')
		remainder = (remainder*10 + digitVal) % 97
	}
	return remainder
}

func validateIBANChecksum(iban string) bool {
	iban = strings.ToUpper(strings.ReplaceAll(iban, " ", ""))

	if len(iban) < 4 {
		return false
	}

	rearranged := iban[4:] + iban[0:4]

	numericString := ""
	for _, char := range rearranged {
		if char >= 'A' && char <= 'Z' {
			numericString += fmt.Sprintf("%d", int(char)-'A'+10)
		} else if char >= '0' && char <= '9' {
			numericString += string(char)
		} else {
			return false
		}
	}

	remainder := calculateMod97(numericString)
	return remainder == 1
}

// ValidateIBAN validates an IBAN with comprehensive checks
func ValidateIBAN(iban string) models.IBANValidation {
	result := models.IBANValidation{
		IBAN:               iban,
		IsValid:            false,
		FormattedIBAN:      "",
		CountryCode:        "",
		CountryName:        "",
		CheckDigits:        "",
		BBAN:               "",
		BankCode:           "",
		AccountNumber:      "",
		IsFormatValid:      false,
		IsCountrySupported: false,
		IsLengthValid:      false,
		IsChecksumValid:    false,
	}

	cleanIBAN := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(iban), " ", ""))

	if len(cleanIBAN) < 15 {
		return result
	}

	if len(cleanIBAN) < 2 || !isIBANLetter(rune(cleanIBAN[0])) || !isIBANLetter(rune(cleanIBAN[1])) {
		return result
	}
	result.CountryCode = cleanIBAN[0:2]

	if len(cleanIBAN) < 4 || !isIBANDigit(rune(cleanIBAN[2])) || !isIBANDigit(rune(cleanIBAN[3])) {
		return result
	}
	result.CheckDigits = cleanIBAN[2:4]

	specs := models.GetIBANCountrySpecs()
	spec, exists := specs[result.CountryCode]
	if !exists {
		return result
	}
	result.IsCountrySupported = true
	result.CountryName = spec.CountryName

	if len(cleanIBAN) != spec.Length {
		return result
	}
	result.IsLengthValid = true

	result.BBAN = cleanIBAN[4:]

	bbanRegex := regexp.MustCompile(spec.BBANFormat)
	if !bbanRegex.MatchString(result.BBAN) {
		return result
	}
	result.IsFormatValid = true

	if spec.BankCodeLen > 0 && spec.BankCodeStart+spec.BankCodeLen <= len(cleanIBAN) {
		result.BankCode = cleanIBAN[spec.BankCodeStart : spec.BankCodeStart+spec.BankCodeLen]
	}
	if spec.AccountLen > 0 && spec.AccountStart+spec.AccountLen <= len(cleanIBAN) {
		result.AccountNumber = cleanIBAN[spec.AccountStart : spec.AccountStart+spec.AccountLen]
	}

	result.IsChecksumValid = validateIBANChecksum(cleanIBAN)
	result.FormattedIBAN = formatIBANWithSpaces(cleanIBAN)

	result.IsValid = result.IsCountrySupported && result.IsLengthValid &&
		result.IsFormatValid && result.IsChecksumValid

	return result
}
EOF
```

**Action:** Copy all IBAN validation functions from `iban-handling.go` (lines 244-end), excluding the structs and country specs map.

### Step 5: Create `internal/services/generator/qr.go`

```bash
# First, read qr_handler.go to extract the business logic
# Then create the service file
cat > internal/services/generator/qr.go << 'EOF'
package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/innovelabs/microtools-go/internal/models"
	qrcode "github.com/skip2/go-qrcode"
)

// Supported QR types
var supportedTypes = map[string]bool{
	"text": true, "url": true, "email": true, "tel": true,
	"sms": true, "wifi": true, "vcard": true, "geo": true,
	"event": true, "json": true,
}

// ApplyDefaults applies default values to QR request
func ApplyDefaults(req *models.QRRequest) {
	if req.Options.Size == 0 {
		req.Options.Size = 256
	}
	if req.Options.ErrorCorrection == "" {
		req.Options.ErrorCorrection = "M"
	}
}

// ValidateRequest validates a QR generation request
func ValidateRequest(req models.QRRequest) error {
	if req.Type == "" {
		return errors.New("type is required")
	}
	if !supportedTypes[req.Type] {
		return fmt.Errorf("unsupported type: %s", req.Type)
	}
	if req.Data == "" && req.Type != "wifi" && req.Type != "vcard" && req.Type != "event" {
		return errors.New("data is required for this type")
	}
	if req.Options.Size < 64 || req.Options.Size > 2048 {
		return errors.New("size must be between 64 and 2048")
	}
	return nil
}

// BuildPayload builds the QR code payload based on type
func BuildPayload(qrType, data string) (string, error) {
	switch qrType {
	case "text":
		return data, nil
	case "url":
		if !strings.HasPrefix(data, "http://") && !strings.HasPrefix(data, "https://") {
			return "", errors.New("URL must start with http:// or https://")
		}
		return data, nil
	case "email":
		return fmt.Sprintf("mailto:%s", data), nil
	case "tel":
		return fmt.Sprintf("tel:%s", data), nil
	case "sms":
		return fmt.Sprintf("sms:%s", data), nil
	case "wifi":
		var wifi models.WifiData
		if err := json.Unmarshal([]byte(data), &wifi); err != nil {
			return "", errors.New("invalid WiFi data format")
		}
		return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", wifi.Security, wifi.SSID, wifi.Password), nil
	case "vcard":
		var vcard models.VCardData
		if err := json.Unmarshal([]byte(data), &vcard); err != nil {
			return "", errors.New("invalid vCard data format")
		}
		return fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s %s\nORG:%s\nTEL:%s\nEMAIL:%s\nEND:VCARD",
			vcard.FirstName, vcard.LastName, vcard.Org, vcard.Phone, vcard.Email), nil
	case "geo":
		return fmt.Sprintf("geo:%s", data), nil
	case "event":
		var event models.EventData
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return "", errors.New("invalid event data format")
		}
		return fmt.Sprintf("BEGIN:VEVENT\nSUMMARY:%s\nDTSTART:%s\nDTEND:%s\nEND:VEVENT",
			event.Summary, event.Start, event.End), nil
	case "json":
		return data, nil
	default:
		return "", fmt.Errorf("unsupported type: %s", qrType)
	}
}

// ParseErrorCorrection parses error correction level
func ParseErrorCorrection(level string) qrcode.RecoveryLevel {
	switch strings.ToUpper(level) {
	case "L":
		return qrcode.Low
	case "M":
		return qrcode.Medium
	case "Q":
		return qrcode.High
	case "H":
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

// GenerateQR generates a QR code PNG image
func GenerateQR(req models.QRRequest) ([]byte, error) {
	ApplyDefaults(&req)

	if err := ValidateRequest(req); err != nil {
		return nil, err
	}

	payload, err := BuildPayload(req.Type, req.Data)
	if err != nil {
		return nil, err
	}

	level := ParseErrorCorrection(req.Options.ErrorCorrection)

	png, err := qrcode.Encode(payload, level, req.Options.Size)
	if err != nil {
		return nil, errors.New("failed to generate QR code")
	}

	return png, nil
}
EOF
```

**Action:** Extract QR generation logic from `qr_handler.go`, separating business logic from HTTP handling.

### Step 6: Create `internal/services/generator/barcode.go`

```bash
cat > internal/services/generator/barcode.go << 'EOF'
package generator

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/ean"
	"github.com/innovelabs/microtools-go/internal/models"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	BarcodeTypeUPCA    = "UPC-A"
	BarcodeTypeEAN13   = "EAN-13"
	BarcodeTypeCode128 = "Code128"

	defaultBarcodeWidth  = 300
	defaultBarcodeHeight = 150
)

var (
	ErrInvalidType   = errors.New("invalid barcode type")
	ErrInvalidFormat = errors.New("invalid format")
	ErrInvalidData   = errors.New("invalid data for barcode type")
)

// BarcodeService defines barcode generation interface
type BarcodeService interface {
	Generate(req models.GenerateRequest) ([]byte, string, error)
}

type defaultBarcodeService struct{}

// NewDefaultBarcodeService creates a new barcode service
func NewDefaultBarcodeService() BarcodeService {
	return &defaultBarcodeService{}
}

// Generate generates a barcode image
func (s *defaultBarcodeService) Generate(req models.GenerateRequest) ([]byte, string, error) {
	// Apply defaults
	if req.Width == 0 {
		req.Width = defaultBarcodeWidth
	}
	if req.Height == 0 {
		req.Height = defaultBarcodeHeight
	}
	if req.Format == "" {
		req.Format = "png"
	}

	// Validate
	if req.Format != "png" && req.Format != "svg" {
		return nil, "", ErrInvalidFormat
	}

	// Generate barcode
	var bc barcode.Barcode
	var err error

	switch req.Type {
	case BarcodeTypeUPCA:
		bc, err = ean.Encode(req.Data)
	case BarcodeTypeEAN13:
		bc, err = ean.Encode(req.Data)
	case BarcodeTypeCode128:
		bc, err = code128.Encode(req.Data)
	default:
		return nil, "", ErrInvalidType
	}

	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrInvalidData, err)
	}

	// Scale barcode
	bc, err = barcode.Scale(bc, req.Width, req.Height)
	if err != nil {
		return nil, "", err
	}

	// Create image
	img := bc

	// Add text if requested
	if req.IncludeText {
		img = addTextToBarcode(bc, req.Data)
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), "image/png", nil
}

func addTextToBarcode(bc barcode.Barcode, text string) image.Image {
	bounds := bc.Bounds()
	textHeight := 20
	newHeight := bounds.Dy() + textHeight

	img := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), newHeight))
	draw.Draw(img, bounds, bc, bounds.Min, draw.Src)

	// Draw text
	point := fixed.Point26_6{
		X: fixed.Int26_6((bounds.Dx() - len(text)*7) / 2 * 64),
		Y: fixed.Int26_6((bounds.Dy() + 15) * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)

	return img
}
EOF
```

**Action:** Extract barcode generation logic from `bar_code_handler.go`.

**Checkpoint Phase 3:**
```bash
go build ./internal/services/...
```

---

## Phase 4: Migrate Handlers (HTTP Layer)

### Step 7: Create `internal/handlers/validation.go`

```bash
cat > internal/handlers/validation.go << 'EOF'
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
EOF
```

### Step 8: Create `internal/handlers/generator.go`

```bash
cat > internal/handlers/generator.go << 'EOF'
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
EOF
```

### Step 9: Create `internal/handlers/user.go`

```bash
cat > internal/handlers/user.go << 'EOF'
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
EOF
```

### Step 10: Create `internal/handlers/health.go`

```bash
cat > internal/handlers/health.go << 'EOF'
package handlers

import (
	"encoding/json"
	"net/http"
)

// LiveHandler handles health check requests
func LiveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Live"})
}
EOF
```

**Checkpoint Phase 4:**
```bash
go build ./internal/handlers
```

---

## Phase 5: Migrate Middleware & Utils

### Step 11: Create `internal/middleware/auth.go`

```bash
cat > internal/middleware/auth.go << 'EOF'
package middleware

import (
	"net/http"
	"strings"

	"github.com/innovelabs/microtools-go/internal/utils"
)

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
EOF
```

### Step 12: Create `internal/middleware/counter.go`

```bash
cat > internal/middleware/counter.go << 'EOF'
package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

var counterNames = map[string]string{
	"/api/v1/validate/email": "email-validate",
	"/api/v1/validate/ip":    "ip-validate",
	"/api/v1/validate/iban":  "iban-validate",
	"/api/v1/generate/qr":    "qr-generate",
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
EOF
```

### Step 13: Create `internal/utils/jwt.go`

```bash
cat > internal/utils/jwt.go << 'EOF'
package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/innovelabs/microtools-go/internal/config"
)

// GenerateJWT generates a JWT token for a user
func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	cfg := config.LoadConfig()
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ValidateJWT validates a JWT token
func ValidateJWT(tokenString string) (string, error) {
	cfg := config.LoadConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["email"].(string), nil
	}
	return "", err
}
EOF
```

### Step 14: Create `internal/database/mongo.go`

```bash
cat > internal/database/mongo.go << 'EOF'
package database

import (
	"context"
	"log"
	"time"

	"github.com/innovelabs/microtools-go/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDB initializes MongoDB client
func InitMongoDB() *mongo.Client {
	cfg := config.LoadConfig()
	log.Println("Initializing MongoDB... with uri: ", cfg.MongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return client
}
EOF
```

### Step 15: Create `internal/database/redis.go`

```bash
cat > internal/database/redis.go << 'EOF'
package database

import (
	"github.com/go-redis/redis/v8"
	"github.com/innovelabs/microtools-go/internal/config"
)

// InitRedis initializes Redis client
func InitRedis() *redis.Client {
	cfg := config.LoadConfig()
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisURI,
	})
}
EOF
```

**Checkpoint Phase 5:**
```bash
go build ./internal/middleware ./internal/utils ./internal/database
```

---

## Phase 6: Create Router

### Step 16: Create `internal/router/router.go`

```bash
cat > internal/router/router.go << 'EOF'
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

	// Web UI routes
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
EOF
```

**Checkpoint Phase 6:**
```bash
go build ./internal/router
```

---

## Phase 7: Create New Main Entry Point

### Step 17: Create `cmd/api/main.go`

```bash
cat > cmd/api/main.go << 'EOF'
package main

import (
	"log"
	"net/http"

	"github.com/innovelabs/microtools-go/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize databases (currently commented out)
	// cfg := config.LoadConfig()
	// handlers.MongoClient = database.InitMongoDB()
	// redisClient := database.InitRedis()

	log.Println("Database initialized")

	// Setup router
	r := router.SetupRouter()

	// Start server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
EOF
```

**Checkpoint Phase 7:**
```bash
go build -o main ./cmd/api
```

---

## Phase 8: Move Templates & Assets

### Step 18: Move templates

```bash
# Create web templates directory if not exists
mkdir -p web/templates/pages

# Move base template
mv views/base.html web/templates/

# Move page templates
mv views/pages/*.html web/templates/pages/

# Verify
ls web/templates/
ls web/templates/pages/
```

### Step 19: Move assets

```bash
# Create assets directory if not exists
mkdir -p assets

# Move GeoIP database
mv geolite-2-city.mmdb assets/

# Verify
ls -lh assets/
```

**Checkpoint Phase 8:**
```bash
# Verify new locations
test -f web/templates/base.html && echo "✓ base.html moved"
test -f web/templates/pages/home.html && echo "✓ home.html moved"
test -f web/templates/pages/email.html && echo "✓ email.html moved"
test -f web/templates/pages/ip.html && echo "✓ ip.html moved"
test -f web/templates/pages/iban.html && echo "✓ iban.html moved"
test -f web/templates/pages/qr.html && echo "✓ qr.html moved"
test -f web/templates/pages/barcode.html && echo "✓ barcode.html moved"
test -f assets/geolite-2-city.mmdb && echo "✓ GeoIP database moved"
```

---

## Phase 9: Test & Verify

### Step 20: Build the application

```bash
cd /home/steinsgate/main/innovelabs/projects/microtools-go/microtools

# Clean previous builds
rm -f main

# Build new version
go build -o main ./cmd/api

# Check for errors
echo $?  # Should output 0 if successful
```

### Step 21: Run the application

```bash
# Start the server
./main
```

**Expected output:**
```
Database initialized
Server started on :8000
```

### Step 22: Test API endpoints

Open a new terminal and test:

```bash
# Test email validation
curl -X POST http://localhost:8000/api/v1/validate/email \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# Test IP validation
curl -X POST http://localhost:8000/api/v1/validate/ip \
  -H "Content-Type: application/json" \
  -d '{"ip":"8.8.8.8"}'

# Test IBAN validation
curl -X POST http://localhost:8000/api/v1/validate/iban \
  -H "Content-Type: application/json" \
  -d '{"iban":"DE89370400440532013000"}'

# Test health endpoint
curl http://localhost:8000/api/v1/live
```

### Step 23: Test web UI

Open browser and test:
- http://localhost:8000/
- http://localhost:8000/email-validation-api
- http://localhost:8000/ip-geolocation-api
- http://localhost:8000/iban-validation-api
- http://localhost:8000/qr-code-generator-api
- http://localhost:8000/barcode-generator-api

---

## Phase 10: Clean Up (Optional)

### Step 24: Remove old files

**⚠️ IMPORTANT: Only do this after thorough testing!**

```bash
# Create backup first
cp -r . ../microtools-backup

# Remove old Go files
rm -f bar_code_handler.go
rm -f config.go
rm -f count-increment-worker.go
rm -f email-handling.go
rm -f geolocation-handling.go
rm -f handlers.go
rm -f iban-handling.go
rm -f middleware.go
rm -f models.go
rm -f qr_handler.go
rm -f redis-client.go
rm -f utils.go
rm -f main.go  # Old main.go

# Remove old views directory
rm -rf views/

# Verify old files are gone
ls *.go 2>/dev/null || echo "✓ Old .go files removed"
```

### Step 25: Update Dockerfile

Edit your `Dockerfile` to use the new build path:

```dockerfile
# Change this line:
# RUN go build -o main .

# To:
RUN go build -o main ./cmd/api
```

### Step 26: Create .env.example

```bash
cat > .env.example << 'EOF'
MONGO_URI=mongodb://localhost:27017
REDIS_URI=localhost:6379
JWT_SECRET=your-secret-key-here
COUNTER_API_KEY=your-counter-api-key
EOF
```

---

## Final Verification Checklist

- [ ] All API endpoints work correctly
- [ ] All web UI pages render properly
- [ ] Application builds without errors
- [ ] No import errors or missing packages
- [ ] Templates load from new `web/templates/` path
- [ ] GeoIP database loads from new `assets/` path
- [ ] Old files removed (after backup)
- [ ] Dockerfile updated
- [ ] .env.example created

---

## Project Structure (Final)

```
microtools-go/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── requests.go
│   │   ├── responses.go
│   │   └── iban.go
│   ├── services/
│   │   ├── validation/
│   │   │   ├── email.go
│   │   │   ├── ip.go
│   │   │   └── iban.go
│   │   └── generator/
│   │       ├── qr.go
│   │       └── barcode.go
│   ├── handlers/
│   │   ├── validation.go
│   │   ├── generator.go
│   │   ├── user.go
│   │   └── health.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── counter.go
│   ├── router/
│   │   └── router.go
│   ├── database/
│   │   ├── mongo.go
│   │   └── redis.go
│   └── utils/
│       └── jwt.go
├── web/
│   └── templates/
│       ├── base.html
│       └── pages/
│           ├── home.html
│           ├── email.html
│           ├── ip.html
│           ├── iban.html
│           ├── qr.html
│           └── barcode.html
├── assets/
│   └── geolite-2-city.mmdb
├── .env
├── .env.example
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

---

## Troubleshooting

### Build errors

**Error:** `package X is not in GOROOT`
**Solution:** Run `go mod tidy` to download missing dependencies

**Error:** `undefined: function/type`
**Solution:** Check import paths - they should use `github.com/innovelabs/microtools-go/internal/...`

### Runtime errors

**Error:** `template: pattern matches no files`
**Solution:** Verify templates are in `web/templates/` directory

**Error:** `no such file or directory: geolite-2-city.mmdb`
**Solution:** Verify GeoIP database is in `assets/` directory and path is updated in code

### Import cycle errors

**Solution:** Ensure clean separation:
- Models should not import services or handlers
- Services should only import models
- Handlers should import models and services
- Router should import handlers and middleware

---

## Next Steps

1. **Add tests** - Create test files for each service
2. **Add documentation** - Document each package with doc.go files
3. **Add CI/CD** - Set up automated testing and deployment
4. **Add logging** - Implement structured logging
5. **Add metrics** - Add Prometheus metrics
6. **Add API docs** - Generate Swagger/OpenAPI documentation

---

## Support

If you encounter issues during migration:
1. Check the troubleshooting section above
2. Verify each checkpoint passed successfully
3. Review error messages carefully
4. Check the plan file at `/home/steinsgate/.claude/plans/swift-swinging-babbage.md`

Good luck with your migration! 🚀
