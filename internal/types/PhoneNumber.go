package types

import (
	"fmt"
	"regexp"
	"strings"
)

type PhoneNumber string

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`) // Разрешает + и 10-15 цифр

func NewPhoneNumber(phone string) (PhoneNumber, error) {
	// Удаляем ВСЕ нецифровые символы (кроме + в начале)
	clean := strings.ReplaceAll(phone, " ", "")
	clean = strings.ReplaceAll(clean, "-", "")
	clean = strings.ReplaceAll(clean, "(", "")
	clean = strings.ReplaceAll(clean, ")", "")

	// Проверяем regex
	if !phoneRegex.MatchString(clean) {
		return "", fmt.Errorf("phone must be 10-15 digits, got: %v", len(clean))
	}

	// Оставляем только цифры (или + цифры)
	if strings.HasPrefix(clean, "+") {
		return PhoneNumber(clean), nil
	}
	return PhoneNumber(clean), nil
}
