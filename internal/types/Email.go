package types

import (
	"errors"
	"net/mail"
	"strings"
)

type Email string

func NewEmail(email string) (Email, error) {
	if !isValidEmail(email) {
		return "", errors.New("invalid email format")
	}
	return Email(strings.ToLower(email)), nil
}

// email validation
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// from custom type to string
func (e Email) String() string {
	return string(e)
}

// GormDataType for gorm
func (e Email) GormDataType() string {
	return "varchar(255)"
}
