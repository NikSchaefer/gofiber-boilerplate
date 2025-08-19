package utils

import (
	"errors"

	"github.com/nyaruka/phonenumbers"
)

func SanitizePhone(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "US")
	if err != nil {
		return "", errors.New("invalid phone number")
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}
