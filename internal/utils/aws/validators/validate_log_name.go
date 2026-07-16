// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
	"unicode"
)

func validateLogName(name string, nameType string, maxLength int, allowedSpecials string) error {
	if len(name) == 0 {
		return fmt.Errorf("%s name cannot be empty", nameType)
	}
	if len(name) > maxLength {
		return fmt.Errorf("%s name cannot exceed %d characters", nameType, maxLength)
	}
	specialChars := []rune(allowedSpecials)
	for _, c := range name {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) && !containsRune(specialChars, c) {
			return fmt.Errorf("%s name contains invalid character: %c", nameType, c)
		}
	}
	return nil
}

func containsRune(chars []rune, target rune) bool {
	for _, c := range chars {
		if c == target {
			return true
		}
	}
	return false
}
