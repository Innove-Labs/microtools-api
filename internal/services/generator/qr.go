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
