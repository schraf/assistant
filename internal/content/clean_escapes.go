package content

import (
	"regexp"
	"strings"
)

// CleanEscapes removes escape characters from the input string.
// It removes common escape sequences such as:
// - Newline (\n)
// - Carriage return (\r)
// - Tab (\t)
// - Backspace (\b)
// - Form feed (\f)
// - Vertical tab (\v)
// - Backslash (\\)
// - Single quote (\')
// - Double quote (\")
// - Bell (\a)
// - And other escape sequences
func CleanEscapes(text string) string {
	result := text

	// Remove escape sequences that use backslash notation first (before removing backslashes)
	// Remove octal escape sequences (\000 to \377)
	octalPattern := regexp.MustCompile("\\\\[0-7]{1,3}")
	result = octalPattern.ReplaceAllString(result, "")

	// Remove hex escape sequences (\x00 to \xFF)
	hexPattern := regexp.MustCompile("\\\\x[0-9A-Fa-f]{1,2}")
	result = hexPattern.ReplaceAllString(result, "")

	// Remove Unicode escape sequences (\u0000 to \uFFFF)
	unicodePattern := regexp.MustCompile("\\\\u[0-9A-Fa-f]{4}")
	result = unicodePattern.ReplaceAllString(result, "")

	// Remove Unicode escape sequences (\U00000000 to \UFFFFFFFF)
	unicodeLongPattern := regexp.MustCompile("\\\\U[0-9A-Fa-f]{8}")
	result = unicodeLongPattern.ReplaceAllString(result, "")

	// Remove common literal escape sequences
	// Newline
	result = strings.ReplaceAll(result, "\n", "")
	// Carriage return
	result = strings.ReplaceAll(result, "\r", "")
	// Tab
	result = strings.ReplaceAll(result, "\t", "")
	// Backspace
	result = strings.ReplaceAll(result, "\b", "")
	// Form feed
	result = strings.ReplaceAll(result, "\f", "")
	// Vertical tab
	result = strings.ReplaceAll(result, "\v", "")
	// Bell/alert
	result = strings.ReplaceAll(result, "\a", "")
	// Backslash (escaped backslash) - remove remaining backslashes
	result = strings.ReplaceAll(result, "\\", "")

	return result
}
