package handlers

import (
	"net/mail"
	"strings"
	"time"
)

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req SignUpRequest) Validate() error {

	// Validation of name field
	if strings.TrimSpace(req.Name) == "" {
		return &ValidationError{Field: "name", Msg: "must not be empty"}
	}

	// Validation of email field
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return &ValidationError{Field: "email", Msg: "must be a valid email address"}
	}
	
	// Validation of password field
	if len(req.Password) < 8 {
		return &ValidationError{Field: "password", Msg: "must be at least 8 characters"}
	}

	return nil
}

type SignUpResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
} 
