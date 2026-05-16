package audit

import (
	"regexp"
	"strings"
)

var (
	secretKeyTerms = []string{
		"password",
		"passwd",
		"passphrase",
		"secret",
		"token",
		"credential",
		"cookie",
		"session",
		"private_key",
		"private-key",
		"apikey",
		"api_key",
	}

	assignmentSecretPattern = regexp.MustCompile(`(?i)\b(password|passwd|passphrase|secret|token|credential|cookie|session|private[_-]?key|api[_-]?key)\b\s*[:=]\s*("[^"]*"|'[^']*'|[^\s,;]+)`)
	bearerTokenPattern      = regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._~+/=-]+`)
	urlUserinfoPattern      = regexp.MustCompile(`://[^/\s:@]+:[^/\s@]+@`)
)

// RedactString removes common inline secret forms from auditable text.
func RedactString(value string) string {
	value = assignmentSecretPattern.ReplaceAllString(value, `$1=`+redacted)
	value = bearerTokenPattern.ReplaceAllString(value, `Bearer `+redacted)
	value = urlUserinfoPattern.ReplaceAllString(value, `://`+redacted+`@`)

	return value
}

// RedactFields returns a copy of fields with secret-looking keys or values scrubbed.
func RedactFields(fields map[string]string) map[string]string {
	clean := make(map[string]string, len(fields))
	for key, value := range fields {
		if isSecretKey(key) {
			clean[key] = redacted
			continue
		}

		clean[key] = RedactString(value)
	}

	return clean
}

func containsInlineSecret(values ...string) bool {
	for _, value := range values {
		if RedactString(value) != value {
			return true
		}
	}

	return false
}

func isSecretKey(key string) bool {
	normalized := strings.ToLower(strings.TrimSpace(key))
	for _, term := range secretKeyTerms {
		if strings.Contains(normalized, term) {
			return true
		}
	}

	return false
}
