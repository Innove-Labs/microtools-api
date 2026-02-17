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

	specs := models.IBANCountrySpecs
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
