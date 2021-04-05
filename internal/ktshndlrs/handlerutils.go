// Package ktshndlrs provides functions that process commands and events.
package ktshndlrs

import (
	"unicode"
)

// isASCII checks if the given string is an ASCII string.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
