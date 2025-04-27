package types

import (
	"errors"
	"regexp"
	"strings"
)

type PhoneNumber string

var phoneNumberRegex = regexp.MustCompile("^[0-9]{10}$")

func NewPhoneNumber(phoneNumber string) (PhoneNumber, error) {
	cleanNumber := strings.ReplaceAll(strings.ReplaceAll(phoneNumber, " ", ""), "-", "")
	if !phoneNumberRegex.MatchString(cleanNumber) {
		return "", errors.New("invalid phone number")
	}
	return PhoneNumber(cleanNumber), nil
}
func (p PhoneNumber) String() string {
	return string(p)
}
