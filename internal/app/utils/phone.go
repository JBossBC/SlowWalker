package utils

import "regexp"

var regex = `^1[3456789]\d{9}$`
var pattern = regexp.MustCompile(regex)

func IsValidPhoneNumber(phoneNumber string) bool {
	return pattern.Match([]byte(phoneNumber))
}
