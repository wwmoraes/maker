package semver

import (
	"fmt"
	"strings"
)

// Constraint is implemented by semantic version rule values that are used
// to evaluate if a given version satisfies it
type Constraint interface {
	fmt.Stringer

	// Match returns true if the target version satisfies this constraint rules
	Match(target Version, includePrerelease bool) bool
	// IsPrerelease returns true if the constraint contains a prerelease version
	IsPrerelease() bool
}

// NewConstraint returns a new semantic version rule with the given
// constraint string
func NewConstraint(constraintStr string) (Constraint, error) {
	// any version
	if isXRange(constraintStr) {
		return NewAny(constraintStr)
	}

	// OR selector
	if strings.Contains(constraintStr, "||") {
		return NewOrGroup(constraintStr)
	}

	// range selector
	if strings.Contains(constraintStr, " - ") {
		return NewRange(constraintStr)
	}

	// AND selector
	if strings.Contains(constraintStr, " ") {
		return NewAndGroup(constraintStr)
	}

	// tilde selector
	if strings.HasPrefix(constraintStr, "~") {
		return NewTilde(constraintStr)
	}

	// caret selector
	if strings.HasPrefix(constraintStr, "^") {
		return NewCaret(constraintStr)
	}

	// more or equals than selector
	if strings.HasPrefix(constraintStr, ">=") {
		return NewGreaterEqual(constraintStr)
	}

	// less or equals than selector
	if strings.HasPrefix(constraintStr, "<=") {
		return NewLessEqual(constraintStr)
	}

	// more than selector
	if strings.HasPrefix(constraintStr, ">") {
		return NewGreaterThan(constraintStr)
	}

	// less than selector
	if strings.HasPrefix(constraintStr, "<") {
		return NewLessThan(constraintStr)
	}

	// explicit exact selector
	if strings.HasPrefix(constraintStr, "=") {
		return NewEqual(constraintStr)
	}

	// plain selector
	return NewEqual(constraintStr)
}
