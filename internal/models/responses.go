package models

// EmailValidation represents the result of email validation
type EmailValidation struct {
	Email          string `json:"email"`
	IsSyntaxValid  bool   `json:"isSyntaxValid"`
	IsDomainValid  bool   `json:"isDomainValid"`
	MxRecordsFound bool   `json:"mxRecordsFound"`
	IsDisposable   bool   `json:"isDisposable"`
}

// GeoIPResponse represents the result of IP geolocation
type GeoIPResponse struct {
	IP        string  `json:"ip"`
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
}

// IBANValidation represents the result of IBAN validation
type IBANValidation struct {
	IBAN               string `json:"iban"`
	IsValid            bool   `json:"isValid"`
	FormattedIBAN      string `json:"formattedIban"`
	CountryCode        string `json:"countryCode"`
	CountryName        string `json:"countryName"`
	CheckDigits        string `json:"checkDigits"`
	BBAN               string `json:"bban"`
	BankCode           string `json:"bankCode"`
	AccountNumber      string `json:"accountNumber"`
	IsFormatValid      bool   `json:"isFormatValid"`
	IsCountrySupported bool   `json:"isCountrySupported"`
	IsLengthValid      bool   `json:"isLengthValid"`
	IsChecksumValid    bool   `json:"isChecksumValid"`
}

// QRErrorResponse represents a QR generation error
type QRErrorResponse struct {
	Error string `json:"error"`
}
