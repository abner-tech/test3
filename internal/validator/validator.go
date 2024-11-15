package validator

import (
	"regexp"
	"slices"
)

// regexp to check if an email is valid
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// check if email provided is valid
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func PermittedValue(value string, permittedValues ...string) bool {
	return slices.Contains(permittedValues, value)
}

// new type named Validator
type Validator struct {
	Errors map[string]string
}

// construct new validator and return a pointer to it
// all validation errors go into thie one Validator instance
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// checking to see if the Validator's map contains any entries
func (v *Validator) IsEmpty() bool {
	return len(v.Errors) == 0
}

// Add a new error entry to the Validator's error map
func (v *Validator) AddError(key string, message string) {
	_, exists := v.Errors[key]
	if !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(acceptable bool, key string, message string) {
	if !acceptable {
		v.AddError(key, message)
	}
}
