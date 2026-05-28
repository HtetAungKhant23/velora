package user

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

var (
	ErrEmailEmpty   = errors.New("email: cannot be empty")
	ErrEmailInvalid = errors.New("email: invalid format")

	emailRE = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))

	if raw == "" {
		return Email{}, ErrEmailEmpty
	}

	if !emailRE.MatchString(raw) {
		return Email{}, ErrEmailInvalid
	}

	return Email{value: raw}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(email Email) bool {
	return e.value == email.value
}
