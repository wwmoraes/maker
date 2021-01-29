package semver

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidIdentifier means an identifier that does not comply with the
	// semantic version specification
	ErrInvalidIdentifier = errors.New("invalid identifier string")
	// ErrInvalidVersion means a version string that does not comply with the
	// semantic version specification
	ErrInvalidVersion = errors.New("invalid version string")
)

// ParseError is returned by Version and Constraint constructors when they fail
// to process the input string
type ParseError struct {
	Func  string
	Input string
	Err   error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("[semver.%s] parsing '%s': %s", e.Func, e.Input, e.Err.Error())
}

func (e *ParseError) Unwrap() error { return e.Err }
