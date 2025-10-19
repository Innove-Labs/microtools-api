package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
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

type EmailValidation struct {
	Email          string `json:"email"`
	IsSyntaxValid  bool   `json:"isSyntaxValid"`
	IsDomainValid  bool   `json:"isDomainValid"`
	MxRecordsFound bool   `json:"mxRecordsFound"`
	IsDisposable   bool   `json:"isDisposable"`
}

func isValidEmailSyntax(email string) bool {
	// Basic email regex pattern
	var emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// Extract domain from email address
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// Domain Name Validation (MX and A Records)
func isValidDomain(domain string) bool {
	// Check if the domain has MX records
	_, err := net.LookupMX(domain)
	if err == nil {
		return true
	}

	// If no MX records, check if it has an A record
	_, err = net.LookupHost(domain)
	return err == nil
}

func verifyMxRecords(email string) bool {
	domain := extractDomain(email)
	if domain == "" {
		fmt.Println("Invalid email format")
		return false
	}

	// Get MX records for the domain
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		fmt.Println("No MX records found for domain", domain)
		return false
	}

	return true
}

// SMTP Verification (optional)
// func smtpVerify(email string) (bool, bool) {

// 	domain := extractDomain(email)
// 	if domain == "" {
// 		fmt.Println("Invalid email format")
// 		return false, false
// 	}

// 	// Get MX records for the domain
// 	mxRecords, err := net.LookupMX(domain)
// 	if err != nil || len(mxRecords) == 0 {
// 		fmt.Println("No MX records found for domain", domain)
// 		return false, false
// 	}

// 	// Attempt to connect to each MX record until successful or all fail
// 	mxRecordsCount := len(mxRecords)
// 	mxRecordFailedDueToPolicies := 0
// 	for _, mx := range mxRecords {
// 		mxHost := mx.Host
// 		addr := fmt.Sprintf("%s:25", mxHost)

// 		// Establish a connection to the SMTP server with a timeout
// 		conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
// 		if err != nil {
// 			fmt.Printf("Failed to connect to SMTP server %s: %v\n", mxHost, err)
// 			continue // Try the next MX server
// 		}
// 		defer conn.Close()

// 		// Start SMTP client
// 		client, err := smtp.NewClient(conn, mxHost)
// 		if err != nil {
// 			fmt.Printf("Failed to initialize SMTP client for %s: %v\n", mxHost, err)
// 			continue
// 		}
// 		defer client.Quit()

// 		// Send HELO command
// 		if err = client.Hello("example.com"); err != nil {
// 			fmt.Printf("SMTP HELO command rejected by %s: %v\n", mxHost, err)
// 			continue
// 		}

// 		// Set a test sender address (this can be any valid email address)
// 		if err = client.Mail("user@example.com"); err != nil {
// 			fmt.Printf("SMTP MAIL FROM command rejected by %s: %v\n", mxHost, err)
// 			continue
// 		}

// 		// Check the recipient address
// 		if err = client.Rcpt(email); err != nil {
// 			// Check if the error message contains specific phrases

// 			if strings.Contains(err.Error(), "550") || strings.Contains(err.Error(), "not allowed") {
// 				fmt.Println("Recipient address does not exist or is not allowed")
// 			}
// 			// Handle other RCPT TO errors (e.g., denied for anti-spam reasons)
// 			fmt.Printf("SMTP RCPT TO command rejected by %s: %v\n", mxHost, err)
// 			mxRecordFailedDueToPolicies++
// 			continue
// 		}

// 		// If we reach here, the email address is valid on this server
// 		return true, true
// 	}

// 	// All MX servers failed or rejected verification attempts
// 	fmt.Println("SMTP verification could not confirm email address validity", mxRecordFailedDueToPolicies, mxRecordsCount)
// 	return false, mxRecordsCount > mxRecordFailedDueToPolicies
// }

// Check for Disposable Email Address (DEA)
func isDisposableEmail(email string) bool {
	// Extract domain from email
	domain := extractDomain(email)
	if domain == "" {
		fmt.Println("Invalid email format")
		return false
	}

	// Check if domain is in the list of disposable email providers
	for _, disposableDomain := range disposableEmailDomains {
		if strings.ToLower(domain) == disposableDomain {
			return true // Disposable email found
		}
	}
	return false // Not a disposable email
}

func ValidateEmail(email string) EmailValidation {
	emailValidationResult := EmailValidation{
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

	// 4. Check for Disposable Email
	if isDisposableEmail(email) {
		emailValidationResult.IsDisposable = true
	}

	return emailValidationResult
}
