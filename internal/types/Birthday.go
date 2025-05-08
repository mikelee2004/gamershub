package types

import (
	"errors"
	"time"
)

const (
	MinAge     = 8
	DateFormat = "2006-01-02 15:04:05"
)

func ValidateBirthday(bdayString string) (time.Time, error) {
	birthday, err := time.Parse(DateFormat, bdayString)
	if err != nil {
		return time.Time{}, errors.New("invalid date format")
	}
	currentTime := time.Now()
	if birthday.After(currentTime) {
		return time.Time{}, errors.New("birthday is in the future. cant do those")
	}

	age := currentTime.Year() - birthday.Year()
	if currentTime.YearDay() < birthday.YearDay() {
		age--
	}
	if age < MinAge {
		return time.Time{}, errors.New("you are too young")
	}
	return birthday, nil
}
