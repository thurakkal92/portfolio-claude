package contact

import (
	"errors"
	"fmt"
)

type Request struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
	Locale    string `json:"locale"`
	Website   string `json:"website"` // honeypot field — must be empty
	IPAddress string `json:"-"`
	UserAgent string `json:"-"`
}

var (
	ErrRateLimited = errors.New("rate limited")
	ErrHoneypot    = errors.New("honeypot tripped")
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", e.Fields)
}

func newValidationError(fields map[string]string) *ValidationError {
	return &ValidationError{Fields: fields}
}
