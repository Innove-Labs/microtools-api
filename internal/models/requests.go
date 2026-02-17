package models

// EmailRequest represents an email validation request
type EmailRequest struct {
	Email string `json:"email"`
}

// IPRequest represents an IP validation/geolocation request
type IPRequest struct {
	IP string `json:"ip"`
}

// IBANRequest represents an IBAN validation request
type IBANRequest struct {
	IBAN string `json:"iban"`
}

// UserRequest represents a user registration request
type UserRequest struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Company string `json:"company"`
	Country string `json:"country"`
}

// QROptions represents QR code generation options
type QROptions struct {
	Size            int    `json:"size"`
	ErrorCorrection string `json:"error_correction"`
}

// QRRequest represents a QR code generation request
type QRRequest struct {
	Type    string    `json:"type"`
	Data    string    `json:"data"`
	Options QROptions `json:"options"`
}

// WifiData represents WiFi QR code data
type WifiData struct {
	SSID     string `json:"ssid"`
	Password string `json:"password"`
	Security string `json:"security"`
}

// VCardData represents vCard QR code data
type VCardData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Org       string `json:"org"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

// EventData represents event QR code data
type EventData struct {
	Summary string `json:"summary"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

// GenerateRequest represents a barcode generation request
type GenerateRequest struct {
	Data              string `json:"data"`
	Type              string `json:"type"`
	Format            string `json:"format"`
	Width             int    `json:"width"`
	Height            int    `json:"height"`
	IncludeText       bool   `json:"include_text"`
	BackgroundColor   string `json:"background_color"`
	ForegroundColor   string `json:"foreground_color"`
	TextColor         string `json:"text_color"`
	TextPosition      string `json:"text_position"`
	FontSize          int    `json:"font_size"`
	Padding           int    `json:"padding"`
}
