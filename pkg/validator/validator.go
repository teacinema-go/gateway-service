package validator

import (
	"regexp"

	"github.com/nyaruka/phonenumbers"
)

func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidE164Phone(phone string) bool {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return false
	}

	if !phonenumbers.IsValidNumber(num) {
		return false
	}

	formatted := phonenumbers.Format(num, phonenumbers.E164)
	return formatted == phone
}
