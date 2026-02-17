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
