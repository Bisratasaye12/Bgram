package infrastructure

import (
	interfaces "BChat/Domain/Interfaces"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// AuthService is a service for handling authentication-related operations.
type AuthService struct{}

// NewEmailPasswordService creates a new instance of AuthService.
func NewEmailPasswordService() interfaces.EmailPasswordInterface{
	return &AuthService{}
}

// ValidateEmail checks if the provided email is in a valid format.
func (s *AuthService) ValidateEmail(email string) error {
	// Regular expression for validating an email address
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidatePassword checks if the provided password meets the security requirements.
func (s *AuthService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}\\|;:'\",.<>/?`~", char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
