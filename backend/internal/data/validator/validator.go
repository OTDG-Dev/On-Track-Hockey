package validator

import (
	"regexp"
	"slices"
)

// map of validation errors
type Validator struct {
	Errors map[string]string
}

// create a new validator
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// valid returns true if errors map doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// adds an error message to the map so long as no entry aslreadt exists for the given key
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// generic function which returns true if a specific value is in a list of permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// matches returns true if a string value matches a specific regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// generic function which returns true if all values in a slice are unique.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
