package utils

import "strings"

func IsValidCompanyEmail(email string) bool {
	return strings.HasSuffix(email, "@remotestate.com")
}

func GetNameFromEmail(email string) string {
	if !strings.HasSuffix(email, "@remotestate.com") {
		return ""
	}

	localPart := strings.Split(email, "@")[0]
	nameParts := strings.Split(localPart, ".")
	if len(nameParts) != 2 {
		return ""
	}

	first := strings.Title(nameParts[0])
	last := strings.Title(nameParts[1])
	return first + " " + last
}
