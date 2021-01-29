package semver

import (
	"fmt"
	"strings"
)

type hyphenRange struct {
	lower, upper Constraint
}

// NewRange matches versions within the lower and upper constraints (inclusive)
func NewRange(constraintStr string) (Constraint, error) {
	rangeParts := strings.SplitN(constraintStr, " - ", 2)
	if len(rangeParts) != 2 || isXRange(rangeParts[0]) || isXRange(rangeParts[1]) {
		return nil, &ParseError{
			Func:  "NewRange",
			Input: constraintStr,
			Err:   ErrInvalidVersion,
		}
	}

	lowerConstraint, err := NewGreaterEqual(rangeParts[0])
	if err != nil {
		return nil, &ParseError{
			Func:  "NewRange",
			Input: constraintStr,
			Err:   err,
		}
	}

	upperConstraint, err := NewLessEqual(rangeParts[1])
	if err != nil {
		return nil, &ParseError{
			Func:  "NewRange",
			Input: constraintStr,
			Err:   err,
		}
	}

	return &hyphenRange{lowerConstraint, upperConstraint}, nil
}

func (source *hyphenRange) Match(version Version, includePrerelease bool) bool {
	includePrerelease = includePrerelease || source.IsPrerelease()

	if !source.lower.Match(version, includePrerelease && !source.lower.IsPrerelease()) {
		return false
	}

	if !source.upper.Match(version, includePrerelease && !source.upper.IsPrerelease()) {
		return false
	}

	return true
}

func (source *hyphenRange) IsPrerelease() bool {
	return source.lower.IsPrerelease() || source.upper.IsPrerelease()
}

func (source *hyphenRange) String() string {
	return fmt.Sprintf("%s - %s", source.lower.String(), source.upper.String())
}
