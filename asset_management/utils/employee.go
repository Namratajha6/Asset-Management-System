package utils

import (
	"net/http"
	"strings"
)

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

func ParseCommaSeparatedParam(r *http.Request, key string) []string {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}
	parts := strings.Split(param, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
