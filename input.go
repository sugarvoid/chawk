// This file contains internal utilities for sanitizing and validating user-provided
// input for data integrity before sending requests to the API.
package chawk

import (
	"fmt"
	"strings"
)

// RequiredString trims leading and trailing whitespace from a string and
// returns an error if the resulting string is empty.
//
// Use this for mandatory fields like usernames and courseIDs.
func RequiredString(value string, fieldName string) (string, error) {
	trimmedInput := strings.TrimSpace(value)
	if trimmedInput == "" {
		return "", fmt.Errorf("%s is required and cannot be empty or only whitespace", fieldName)
	}
	return trimmedInput, nil
}

// OptionalString trims whitespace from a string. Unlike RequiredString,
// it returns the trimmed value even if it is empty, without an error.
//
// Use this for non-mandatory fields where a blank value is acceptable like user's email
func OptionalString(value string) string {
	return strings.TrimSpace(value)
}

// Ptr is a helper that returns a pointer to the value passed in.
// This is particularly useful for filling out "omitempty" pointer fields in
// update structs without needing to declare local variables.
func ToPtr[T any](v T) *T {
	return &v
}
