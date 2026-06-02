package contact

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

const (
	minNameLen    = 2
	maxNameLen    = 120
	minMessageLen = 10
	maxMessageLen = 5000
	maxEmailLen   = 254
)

func validate(req *Request) error {
	fields := map[string]string{}

	name := strings.TrimSpace(req.Name)
	switch {
	case name == "":
		fields["name"] = "Name is required"
	case utf8.RuneCountInString(name) < minNameLen:
		fields["name"] = "Name is too short"
	case utf8.RuneCountInString(name) > maxNameLen:
		fields["name"] = "Name is too long"
	}
	req.Name = name

	email := strings.TrimSpace(req.Email)
	switch {
	case email == "":
		fields["email"] = "Email is required"
	case len(email) > maxEmailLen:
		fields["email"] = "Email is too long"
	default:
		if _, err := mail.ParseAddress(email); err != nil {
			fields["email"] = "Email is not valid"
		}
	}
	req.Email = email

	msg := strings.TrimSpace(req.Message)
	switch {
	case msg == "":
		fields["message"] = "Message is required"
	case utf8.RuneCountInString(msg) < minMessageLen:
		fields["message"] = "Message is too short (min 10 characters)"
	case utf8.RuneCountInString(msg) > maxMessageLen:
		fields["message"] = "Message is too long"
	}
	req.Message = msg

	if len(fields) > 0 {
		return newValidationError(fields)
	}
	return nil
}
