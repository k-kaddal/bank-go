package util

import "regexp"

func IsValidAlphanumeric(str string) bool {
	alphanumeric := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	return alphanumeric.MatchString(str)
}