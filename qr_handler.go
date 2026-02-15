package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

// --- Request / Response structs ---

type QROptions struct {
	Size            int    `json:"size"`
	ErrorCorrection string `json:"error_correction"`
}

type QRRequest struct {
	Type    string    `json:"type"`
	Data    string    `json:"data"`
	Options QROptions `json:"options"`
}

type QRErrorResponse struct {
	Error string `json:"error"`
}

// --- Structs for structured type payloads ---

type WifiData struct {
	SSID     string `json:"ssid"`
	Password string `json:"password"`
	Security string `json:"security"`
}

type VCardData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Org       string `json:"org"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type EventData struct {
	Summary string `json:"summary"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

// Supported QR types. Add new entries here to extend.
var supportedTypes = map[string]bool{
	"text": true, "url": true, "email": true, "tel": true,
	"sms": true, "wifi": true, "vcard": true, "geo": true,
	"event": true, "json": true,
}

// --- HTTP handler ---

// QRHandler handles POST /generate requests.
// It decodes the JSON body, validates inputs, formats the payload
// according to the requested type, and returns a PNG QR code image.
func QRHandler(w http.ResponseWriter, r *http.Request) {
	var req QRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Apply defaults before validation.
	applyDefaults(&req)

	if err := validateRequest(req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	payload, err := buildPayloadByType(req.Type, req.Data)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	level := parseErrorCorrection(req.Options.ErrorCorrection)

	png, err := qrcode.Encode(payload, level, req.Options.Size)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to generate QR code")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(png)
}

// --- Defaults ---

func applyDefaults(req *QRRequest) {
	if req.Type == "" {
		req.Type = "text"
	}
	if req.Options.Size == 0 {
		req.Options.Size = 256
	}
	if req.Options.ErrorCorrection == "" {
		req.Options.ErrorCorrection = "medium"
	}
}

// --- Validation ---

// validateRequest checks that all fields meet the documented constraints.
func validateRequest(req QRRequest) error {
	if req.Data == "" {
		return fmt.Errorf("data is required")
	}
	if len(req.Data) > 2000 {
		return fmt.Errorf("data exceeds maximum length of 2000 characters")
	}
	if !supportedTypes[req.Type] {
		return fmt.Errorf("unsupported type: %s", req.Type)
	}
	if req.Options.Size < 128 || req.Options.Size > 1024 {
		return fmt.Errorf("size must be between 128 and 1024")
	}
	ecLevels := map[string]bool{"low": true, "medium": true, "high": true, "highest": true}
	if !ecLevels[req.Options.ErrorCorrection] {
		return fmt.Errorf("error_correction must be one of: low, medium, high, highest")
	}
	return nil
}

// --- Type formatting ---

// buildPayloadByType converts the raw data string into the appropriate
// QR code payload format based on the requested type.
func buildPayloadByType(qrType, data string) (string, error) {
	switch qrType {
	case "text":
		return data, nil

	case "url":
		if !strings.HasPrefix(data, "http://") && !strings.HasPrefix(data, "https://") {
			return "", fmt.Errorf("url must start with http:// or https://")
		}
		return data, nil

	case "email":
		return "mailto:" + data, nil

	case "tel":
		return "tel:" + data, nil

	case "sms":
		return buildSMS(data)

	case "wifi":
		return buildWifi(data)

	case "vcard":
		return buildVCard(data)

	case "geo":
		return buildGeo(data)

	case "event":
		return buildEvent(data)

	case "json":
		if !json.Valid([]byte(data)) {
			return "", fmt.Errorf("data is not valid JSON")
		}
		return data, nil

	default:
		return "", fmt.Errorf("unsupported type: %s", qrType)
	}
}

// buildSMS expects "number|message" and formats as smsto:number:message.
func buildSMS(data string) (string, error) {
	parts := strings.SplitN(data, "|", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("sms data must be in format: number|message")
	}
	return fmt.Sprintf("smsto:%s:%s", parts[0], parts[1]), nil
}

// buildWifi parses a JSON string with ssid, password, security and
// produces the standard WIFI: QR format.
func buildWifi(data string) (string, error) {
	var w WifiData
	if err := json.Unmarshal([]byte(data), &w); err != nil {
		return "", fmt.Errorf("wifi data must be valid JSON with ssid, password, security")
	}
	if w.SSID == "" {
		return "", fmt.Errorf("wifi ssid is required")
	}
	validSec := map[string]bool{"WPA": true, "WEP": true, "nopass": true}
	if !validSec[w.Security] {
		return "", fmt.Errorf("wifi security must be one of: WPA, WEP, nopass")
	}
	return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", w.Security, w.SSID, w.Password), nil
}

// buildVCard parses JSON with contact fields and builds a VCARD 3.0 string.
func buildVCard(data string) (string, error) {
	var v VCardData
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return "", fmt.Errorf("vcard data must be valid JSON with first_name, last_name, org, phone, email")
	}
	if v.FirstName == "" || v.LastName == "" {
		return "", fmt.Errorf("vcard requires first_name and last_name")
	}

	var b strings.Builder
	b.WriteString("BEGIN:VCARD\r\n")
	b.WriteString("VERSION:3.0\r\n")
	b.WriteString(fmt.Sprintf("N:%s;%s;;;\r\n", v.LastName, v.FirstName))
	b.WriteString(fmt.Sprintf("FN:%s %s\r\n", v.FirstName, v.LastName))
	if v.Org != "" {
		b.WriteString(fmt.Sprintf("ORG:%s\r\n", v.Org))
	}
	if v.Phone != "" {
		b.WriteString(fmt.Sprintf("TEL:%s\r\n", v.Phone))
	}
	if v.Email != "" {
		b.WriteString(fmt.Sprintf("EMAIL:%s\r\n", v.Email))
	}
	b.WriteString("END:VCARD")
	return b.String(), nil
}

// buildGeo expects "lat,long" and formats as geo:lat,long.
func buildGeo(data string) (string, error) {
	parts := strings.SplitN(data, ",", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return "", fmt.Errorf("geo data must be in format: lat,long")
	}
	return fmt.Sprintf("geo:%s,%s", strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])), nil
}

// buildEvent parses JSON with summary, start, end (ISO 8601) and builds
// a minimal VCALENDAR/VEVENT string.
func buildEvent(data string) (string, error) {
	var e EventData
	if err := json.Unmarshal([]byte(data), &e); err != nil {
		return "", fmt.Errorf("event data must be valid JSON with summary, start, end")
	}
	if e.Summary == "" || e.Start == "" || e.End == "" {
		return "", fmt.Errorf("event requires summary, start, and end")
	}

	startTime, err := time.Parse(time.RFC3339, e.Start)
	if err != nil {
		return "", fmt.Errorf("start must be a valid ISO 8601 datetime")
	}
	endTime, err := time.Parse(time.RFC3339, e.End)
	if err != nil {
		return "", fmt.Errorf("end must be a valid ISO 8601 datetime")
	}

	const vcalFmt = "20060102T150405Z"
	var b strings.Builder
	b.WriteString("BEGIN:VCALENDAR\r\n")
	b.WriteString("VERSION:2.0\r\n")
	b.WriteString("BEGIN:VEVENT\r\n")
	b.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", e.Summary))
	b.WriteString(fmt.Sprintf("DTSTART:%s\r\n", startTime.UTC().Format(vcalFmt)))
	b.WriteString(fmt.Sprintf("DTEND:%s\r\n", endTime.UTC().Format(vcalFmt)))
	b.WriteString("END:VEVENT\r\n")
	b.WriteString("END:VCALENDAR")
	return b.String(), nil
}

// --- Error correction mapping ---

// parseErrorCorrection maps the string level name to the qrcode library constant.
func parseErrorCorrection(level string) qrcode.RecoveryLevel {
	switch level {
	case "low":
		return qrcode.Low
	case "medium":
		return qrcode.Medium
	case "high":
		return qrcode.High
	case "highest":
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

// --- Helper ---

// writeJSONError sends a JSON-formatted error response.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(QRErrorResponse{Error: message})
}
