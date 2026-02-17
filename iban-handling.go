package main

import (
	"fmt"
	"regexp"
	"strings"
)

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

// ibanCountrySpecs contains IBAN specifications for 60+ countries
var ibanCountrySpecs = map[string]IBANCountrySpec{
	// SEPA Countries (European Union)
	"AD": {CountryCode: "AD", CountryName: "Andorra", Length: 24, BBANFormat: "^[0-9]{8}[A-Z0-9]{12}$",
		BankCodeStart: 4, BankCodeLen: 8, AccountStart: 12, AccountLen: 12,
		Example: "AD1200012030200359100100"},
	"AT": {CountryCode: "AT", CountryName: "Austria", Length: 20, BBANFormat: "^[0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 5, AccountStart: 9, AccountLen: 11,
		Example: "AT611904300234573201"},
	"BE": {CountryCode: "BE", CountryName: "Belgium", Length: 16, BBANFormat: "^[0-9]{12}$",
		BankCodeStart: 4, BankCodeLen: 3, AccountStart: 7, AccountLen: 9,
		Example: "BE68539007547034"},
	"BG": {CountryCode: "BG", CountryName: "Bulgaria", Length: 22, BBANFormat: "^[A-Z]{4}[0-9]{6}[A-Z0-9]{8}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 14,
		Example: "BG80BNBG96611020345678"},
	"CH": {CountryCode: "CH", CountryName: "Switzerland", Length: 21, BBANFormat: "^[0-9]{5}[A-Z0-9]{12}$",
		BankCodeStart: 4, BankCodeLen: 5, AccountStart: 9, AccountLen: 12,
		Example: "CH9300762011623852957"},
	"CY": {CountryCode: "CY", CountryName: "Cyprus", Length: 28, BBANFormat: "^[0-9]{8}[A-Z0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 8, AccountStart: 12, AccountLen: 16,
		Example: "CY17002001280000001200527600"},
	"CZ": {CountryCode: "CZ", CountryName: "Czech Republic", Length: 24, BBANFormat: "^[0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 16,
		Example: "CZ6508000000192000145399"},
	"DE": {CountryCode: "DE", CountryName: "Germany", Length: 22, BBANFormat: "^[0-9]{18}$",
		BankCodeStart: 4, BankCodeLen: 8, AccountStart: 12, AccountLen: 10,
		Example: "DE89370400440532013000"},
	"DK": {CountryCode: "DK", CountryName: "Denmark", Length: 18, BBANFormat: "^[0-9]{14}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 10,
		Example: "DK5000400440116243"},
	"EE": {CountryCode: "EE", CountryName: "Estonia", Length: 20, BBANFormat: "^[0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 2, AccountStart: 6, AccountLen: 14,
		Example: "EE382200221020145685"},
	"ES": {CountryCode: "ES", CountryName: "Spain", Length: 24, BBANFormat: "^[0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 8, AccountStart: 12, AccountLen: 12,
		Example: "ES9121000418450200051332"},
	"FI": {CountryCode: "FI", CountryName: "Finland", Length: 18, BBANFormat: "^[0-9]{14}$",
		BankCodeStart: 4, BankCodeLen: 6, AccountStart: 10, AccountLen: 8,
		Example: "FI2112345600000785"},
	"FR": {CountryCode: "FR", CountryName: "France", Length: 27, BBANFormat: "^[0-9]{10}[A-Z0-9]{11}[0-9]{2}$",
		BankCodeStart: 4, BankCodeLen: 10, AccountStart: 14, AccountLen: 13,
		Example: "FR1420041010050500013M02606"},
	"GB": {CountryCode: "GB", CountryName: "United Kingdom", Length: 22, BBANFormat: "^[A-Z]{4}[0-9]{14}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 14,
		Example: "GB29NWBK60161331926819"},
	"GI": {CountryCode: "GI", CountryName: "Gibraltar", Length: 23, BBANFormat: "^[A-Z]{4}[A-Z0-9]{15}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 15,
		Example: "GI75NWBK000000007099453"},
	"GR": {CountryCode: "GR", CountryName: "Greece", Length: 27, BBANFormat: "^[0-9]{7}[A-Z0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 7, AccountStart: 11, AccountLen: 16,
		Example: "GR1601101250000000012300695"},
	"HR": {CountryCode: "HR", CountryName: "Croatia", Length: 21, BBANFormat: "^[0-9]{17}$",
		BankCodeStart: 4, BankCodeLen: 7, AccountStart: 11, AccountLen: 10,
		Example: "HR1210010051863000160"},
	"HU": {CountryCode: "HU", CountryName: "Hungary", Length: 28, BBANFormat: "^[0-9]{24}$",
		BankCodeStart: 4, BankCodeLen: 7, AccountStart: 11, AccountLen: 17,
		Example: "HU42117730161111101800000000"},
	"IE": {CountryCode: "IE", CountryName: "Ireland", Length: 22, BBANFormat: "^[A-Z]{4}[0-9]{14}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 14,
		Example: "IE29AIBK93115212345678"},
	"IS": {CountryCode: "IS", CountryName: "Iceland", Length: 26, BBANFormat: "^[0-9]{22}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 18,
		Example: "IS140159260076545510730339"},
	"IT": {CountryCode: "IT", CountryName: "Italy", Length: 27, BBANFormat: "^[A-Z][0-9]{10}[A-Z0-9]{12}$",
		BankCodeStart: 5, BankCodeLen: 10, AccountStart: 15, AccountLen: 12,
		Example: "IT60X0542811101000000123456"},
	"LI": {CountryCode: "LI", CountryName: "Liechtenstein", Length: 21, BBANFormat: "^[0-9]{5}[A-Z0-9]{12}$",
		BankCodeStart: 4, BankCodeLen: 5, AccountStart: 9, AccountLen: 12,
		Example: "LI21088100002324013AA"},
	"LT": {CountryCode: "LT", CountryName: "Lithuania", Length: 20, BBANFormat: "^[0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 5, AccountStart: 9, AccountLen: 11,
		Example: "LT121000011101001000"},
	"LU": {CountryCode: "LU", CountryName: "Luxembourg", Length: 20, BBANFormat: "^[0-9]{3}[A-Z0-9]{13}$",
		BankCodeStart: 4, BankCodeLen: 3, AccountStart: 7, AccountLen: 13,
		Example: "LU280019400644750000"},
	"LV": {CountryCode: "LV", CountryName: "Latvia", Length: 21, BBANFormat: "^[A-Z]{4}[A-Z0-9]{13}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 13,
		Example: "LV80BANK0000435195001"},
	"MC": {CountryCode: "MC", CountryName: "Monaco", Length: 27, BBANFormat: "^[0-9]{10}[A-Z0-9]{11}[0-9]{2}$",
		BankCodeStart: 4, BankCodeLen: 10, AccountStart: 14, AccountLen: 13,
		Example: "MC5811222000010123456789030"},
	"MT": {CountryCode: "MT", CountryName: "Malta", Length: 31, BBANFormat: "^[A-Z]{4}[0-9]{5}[A-Z0-9]{18}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 23,
		Example: "MT84MALT011000012345MTLCAST001S"},
	"NL": {CountryCode: "NL", CountryName: "Netherlands", Length: 18, BBANFormat: "^[A-Z]{4}[0-9]{10}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 10,
		Example: "NL91ABNA0417164300"},
	"NO": {CountryCode: "NO", CountryName: "Norway", Length: 15, BBANFormat: "^[0-9]{11}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 7,
		Example: "NO9386011117947"},
	"PL": {CountryCode: "PL", CountryName: "Poland", Length: 28, BBANFormat: "^[0-9]{24}$",
		BankCodeStart: 4, BankCodeLen: 8, AccountStart: 12, AccountLen: 16,
		Example: "PL61109010140000071219812874"},
	"PT": {CountryCode: "PT", CountryName: "Portugal", Length: 25, BBANFormat: "^[0-9]{21}$",
		BankCodeStart: 4, BankCodeLen: 8, AccountStart: 12, AccountLen: 13,
		Example: "PT50000201231234567890154"},
	"RO": {CountryCode: "RO", CountryName: "Romania", Length: 24, BBANFormat: "^[A-Z]{4}[A-Z0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 16,
		Example: "RO49AAAA1B31007593840000"},
	"SE": {CountryCode: "SE", CountryName: "Sweden", Length: 24, BBANFormat: "^[0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 3, AccountStart: 7, AccountLen: 17,
		Example: "SE4550000000058398257466"},
	"SI": {CountryCode: "SI", CountryName: "Slovenia", Length: 19, BBANFormat: "^[0-9]{15}$",
		BankCodeStart: 4, BankCodeLen: 5, AccountStart: 9, AccountLen: 10,
		Example: "SI56263300012039086"},
	"SK": {CountryCode: "SK", CountryName: "Slovakia", Length: 24, BBANFormat: "^[0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 16,
		Example: "SK3112000000198742637541"},

	// Non-SEPA European Countries
	"SM": {CountryCode: "SM", CountryName: "San Marino", Length: 27, BBANFormat: "^[A-Z][0-9]{10}[A-Z0-9]{12}$",
		BankCodeStart: 5, BankCodeLen: 10, AccountStart: 15, AccountLen: 12,
		Example: "SM86U0322509800000000270100"},
	"VA": {CountryCode: "VA", CountryName: "Vatican City", Length: 22, BBANFormat: "^[0-9]{18}$",
		BankCodeStart: 4, BankCodeLen: 3, AccountStart: 7, AccountLen: 15,
		Example: "VA59001123000012345678"},

	// Middle East & North Africa
	"AE": {CountryCode: "AE", CountryName: "United Arab Emirates", Length: 23, BBANFormat: "^[0-9]{19}$",
		BankCodeStart: 4, BankCodeLen: 3, AccountStart: 7, AccountLen: 16,
		Example: "AE070331234567890123456"},
	"BH": {CountryCode: "BH", CountryName: "Bahrain", Length: 22, BBANFormat: "^[A-Z]{4}[A-Z0-9]{14}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 14,
		Example: "BH67BMAG00001299123456"},
	"IL": {CountryCode: "IL", CountryName: "Israel", Length: 23, BBANFormat: "^[0-9]{19}$",
		BankCodeStart: 4, BankCodeLen: 6, AccountStart: 10, AccountLen: 13,
		Example: "IL620108000000099999999"},
	"JO": {CountryCode: "JO", CountryName: "Jordan", Length: 30, BBANFormat: "^[A-Z]{4}[0-9]{4}[A-Z0-9]{18}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 22,
		Example: "JO94CBJO0010000000000131000302"},
	"KW": {CountryCode: "KW", CountryName: "Kuwait", Length: 30, BBANFormat: "^[A-Z]{4}[A-Z0-9]{22}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 22,
		Example: "KW81CBKU0000000000001234560101"},
	"LB": {CountryCode: "LB", CountryName: "Lebanon", Length: 28, BBANFormat: "^[0-9]{4}[A-Z0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 20,
		Example: "LB62099900000001001901229114"},
	"PS": {CountryCode: "PS", CountryName: "Palestine", Length: 29, BBANFormat: "^[A-Z]{4}[A-Z0-9]{21}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 21,
		Example: "PS92PALS000000000400123456702"},
	"QA": {CountryCode: "QA", CountryName: "Qatar", Length: 29, BBANFormat: "^[A-Z]{4}[A-Z0-9]{21}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 21,
		Example: "QA58DOHB00001234567890ABCDEFG"},
	"SA": {CountryCode: "SA", CountryName: "Saudi Arabia", Length: 24, BBANFormat: "^[0-9]{2}[A-Z0-9]{18}$",
		BankCodeStart: 4, BankCodeLen: 2, AccountStart: 6, AccountLen: 18,
		Example: "SA0380000000608010167519"},
	"TR": {CountryCode: "TR", CountryName: "Turkey", Length: 26, BBANFormat: "^[0-9]{5}[A-Z0-9]{17}$",
		BankCodeStart: 4, BankCodeLen: 5, AccountStart: 9, AccountLen: 17,
		Example: "TR330006100519786457841326"},

	// Caribbean & Latin America
	"BR": {CountryCode: "BR", CountryName: "Brazil", Length: 29, BBANFormat: "^[0-9]{23}[A-Z][A-Z0-9]$",
		BankCodeStart: 4, BankCodeLen: 8, AccountStart: 12, AccountLen: 17,
		Example: "BR1800360305000010009795493C1"},
	"CR": {CountryCode: "CR", CountryName: "Costa Rica", Length: 22, BBANFormat: "^[0-9]{18}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 14,
		Example: "CR05015202001026284066"},
	"DO": {CountryCode: "DO", CountryName: "Dominican Republic", Length: 28, BBANFormat: "^[A-Z]{4}[0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 20,
		Example: "DO28BAGR00000001212453611324"},
	"GT": {CountryCode: "GT", CountryName: "Guatemala", Length: 28, BBANFormat: "^[A-Z0-9]{24}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 20,
		Example: "GT82TRAJ01020000001210029690"},
	"SV": {CountryCode: "SV", CountryName: "El Salvador", Length: 28, BBANFormat: "^[A-Z]{4}[0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 20,
		Example: "SV62CENR00000000000000700025"},

	// Other regions
	"AZ": {CountryCode: "AZ", CountryName: "Azerbaijan", Length: 28, BBANFormat: "^[A-Z]{4}[A-Z0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 20,
		Example: "AZ21NABZ00000000137010001944"},
	"BY": {CountryCode: "BY", CountryName: "Belarus", Length: 28, BBANFormat: "^[A-Z0-9]{4}[0-9]{4}[A-Z0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 20,
		Example: "BY13NBRB3600900000002Z00AB00"},
	"EG": {CountryCode: "EG", CountryName: "Egypt", Length: 29, BBANFormat: "^[0-9]{25}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 21,
		Example: "EG380019000500000000263180002"},
	"GE": {CountryCode: "GE", CountryName: "Georgia", Length: 22, BBANFormat: "^[A-Z]{2}[0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 2, AccountStart: 6, AccountLen: 16,
		Example: "GE29NB0000000101904917"},
	"IQ": {CountryCode: "IQ", CountryName: "Iraq", Length: 23, BBANFormat: "^[A-Z]{4}[0-9]{15}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 15,
		Example: "IQ98NBIQ850123456789012"},
	"KZ": {CountryCode: "KZ", CountryName: "Kazakhstan", Length: 20, BBANFormat: "^[0-9]{3}[A-Z0-9]{13}$",
		BankCodeStart: 4, BankCodeLen: 3, AccountStart: 7, AccountLen: 13,
		Example: "KZ86125KZT5004100100"},
	"MD": {CountryCode: "MD", CountryName: "Moldova", Length: 24, BBANFormat: "^[A-Z0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 2, AccountStart: 6, AccountLen: 18,
		Example: "MD24AG000225100013104168"},
	"MU": {CountryCode: "MU", CountryName: "Mauritius", Length: 30, BBANFormat: "^[A-Z]{4}[0-9]{19}[A-Z]{3}$",
		BankCodeStart: 4, BankCodeLen: 6, AccountStart: 10, AccountLen: 20,
		Example: "MU17BOMM0101101030300200000MUR"},
	"PK": {CountryCode: "PK", CountryName: "Pakistan", Length: 24, BBANFormat: "^[A-Z]{4}[A-Z0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 16,
		Example: "PK36SCBL0000001123456702"},
	"TN": {CountryCode: "TN", CountryName: "Tunisia", Length: 24, BBANFormat: "^[0-9]{20}$",
		BankCodeStart: 4, BankCodeLen: 5, AccountStart: 9, AccountLen: 15,
		Example: "TN5910006035183598478831"},
	"UA": {CountryCode: "UA", CountryName: "Ukraine", Length: 29, BBANFormat: "^[0-9]{6}[A-Z0-9]{19}$",
		BankCodeStart: 4, BankCodeLen: 6, AccountStart: 10, AccountLen: 19,
		Example: "UA213223130000026007233566001"},
	"XK": {CountryCode: "XK", CountryName: "Kosovo", Length: 20, BBANFormat: "^[0-9]{16}$",
		BankCodeStart: 4, BankCodeLen: 4, AccountStart: 8, AccountLen: 12,
		Example: "XK051212012345678906"},
}

// Helper function to check if a rune is a letter
func isIBANLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

// Helper function to check if a rune is a digit
func isIBANDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// Helper function to format IBAN with spaces every 4 characters
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

// calculateMod97 calculates mod 97 of a very large number represented as a string
func calculateMod97(numStr string) int {
	remainder := 0
	for _, digit := range numStr {
		digitVal := int(digit - '0')
		remainder = (remainder*10 + digitVal) % 97
	}
	return remainder
}

// validateIBANChecksum validates the IBAN checksum using the mod-97 algorithm (ISO 13616)
func validateIBANChecksum(iban string) bool {
	// Remove spaces and convert to uppercase
	iban = strings.ToUpper(strings.ReplaceAll(iban, " ", ""))

	if len(iban) < 4 {
		return false
	}

	// Step 1: Move first 4 characters to end
	rearranged := iban[4:] + iban[0:4]

	// Step 2: Replace letters with numbers (A=10, B=11, ..., Z=35)
	numericString := ""
	for _, char := range rearranged {
		if char >= 'A' && char <= 'Z' {
			// A=10, B=11, etc.
			numericString += fmt.Sprintf("%d", int(char)-'A'+10)
		} else if char >= '0' && char <= '9' {
			numericString += string(char)
		} else {
			return false // Invalid character
		}
	}

	// Step 3: Calculate mod 97
	remainder := calculateMod97(numericString)

	// Step 4: Valid if remainder equals 1
	return remainder == 1
}

// ValidateIBAN validates an IBAN with comprehensive checks
func ValidateIBAN(iban string) IBANValidation {
	result := IBANValidation{
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

	// Clean input: remove spaces, convert to uppercase
	cleanIBAN := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(iban), " ", ""))

	// Step 1: Basic format validation (at least 15 chars, starts with 2 letters + 2 digits)
	if len(cleanIBAN) < 15 {
		return result
	}

	// Extract country code (first 2 chars)
	if len(cleanIBAN) < 2 || !isIBANLetter(rune(cleanIBAN[0])) || !isIBANLetter(rune(cleanIBAN[1])) {
		return result
	}
	result.CountryCode = cleanIBAN[0:2]

	// Extract check digits (chars 3-4)
	if len(cleanIBAN) < 4 || !isIBANDigit(rune(cleanIBAN[2])) || !isIBANDigit(rune(cleanIBAN[3])) {
		return result
	}
	result.CheckDigits = cleanIBAN[2:4]

	// Step 2: Check if country is supported
	spec, exists := ibanCountrySpecs[result.CountryCode]
	if !exists {
		return result
	}
	result.IsCountrySupported = true
	result.CountryName = spec.CountryName

	// Step 3: Validate length
	if len(cleanIBAN) != spec.Length {
		return result
	}
	result.IsLengthValid = true

	// Extract BBAN (everything after country code and check digits)
	result.BBAN = cleanIBAN[4:]

	// Step 4: Validate BBAN format using regex
	bbanRegex := regexp.MustCompile(spec.BBANFormat)
	if !bbanRegex.MatchString(result.BBAN) {
		return result
	}
	result.IsFormatValid = true

	// Step 5: Extract bank code and account number
	if spec.BankCodeLen > 0 && spec.BankCodeStart+spec.BankCodeLen <= len(cleanIBAN) {
		result.BankCode = cleanIBAN[spec.BankCodeStart : spec.BankCodeStart+spec.BankCodeLen]
	}
	if spec.AccountLen > 0 && spec.AccountStart+spec.AccountLen <= len(cleanIBAN) {
		result.AccountNumber = cleanIBAN[spec.AccountStart : spec.AccountStart+spec.AccountLen]
	}

	// Step 6: Mod-97 checksum validation
	result.IsChecksumValid = validateIBANChecksum(cleanIBAN)

	// Step 7: Format IBAN with spaces (every 4 characters)
	result.FormattedIBAN = formatIBANWithSpaces(cleanIBAN)

	// Overall validation
	result.IsValid = result.IsCountrySupported && result.IsLengthValid &&
		result.IsFormatValid && result.IsChecksumValid

	return result
}
