package data

import (
	"errors"
)

// custom record error message
var ErrRecordNotFound = errors.New("record not found")

// custom email error message
var ErrDuplicateEmail = errors.New("duplicate email encountered")

var ErrEditConfilct = errors.New("edit confict")
