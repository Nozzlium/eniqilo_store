package util

import (
	"regexp"
)

func ValidatePhoneNumber(phoneNumber string) bool {
	phoneNumberRegexString := `^[+]{1}(?:[0-9-\(\)\/.]\s?){6,15}[0-9]{1}$`
	phoneNumberRegex := regexp.MustCompile(phoneNumberRegexString)

	return phoneNumberRegex.MatchString(phoneNumber)
}
