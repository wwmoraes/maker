package semver

import (
	"fmt"
	"strings"
)

type tilde struct {
	version PartialVersion
}

// NewTilde creates a constraint that allows patch versions within a minor, if
// one was specified, or within a major otherwise
func NewTilde(constraintStr string) (Constraint, error) {
	saneVersionStr := strings.Trim(strings.TrimPrefix(constraintStr, "~"), " ")
	if len(saneVersionStr) == 0 || isXRange(saneVersionStr) {
		return nil, &ParseError{
			Func:  "NewTilde",
			Input: constraintStr,
			Err:   ErrInvalidVersion,
		}
	}

	version, err := NewPartialVersion(saneVersionStr)
	if err != nil {
		return nil, &ParseError{
			Func:  "NewTilde",
			Input: constraintStr,
			Err:   err,
		}
	}

	return NewTildeWith(version), nil
}

// NewTildeWith creates a constraint that allows patch versions within a minor, if
// one was specified, or within a major otherwise
func NewTildeWith(version PartialVersion) Constraint {
	return &tilde{
		version,
	}
}

func (source *tilde) Match(target Version, includePrerelease bool) bool {
	if target.IsPrerelease() {
		// do not match prereleases
		if !source.IsPrerelease() {
			return false
		}
	}

	// anything less than the constraint version is not a match
	fullDiff := source.version.Compare(target)

	if fullDiff == -1 {
		return false
	}

	// major didn't match
	if source.version.Major() >= 0 && source.version.CompareMajor(target) != 0 {
		return false
	}

	// minor didn't match
	if source.version.Minor() >= 0 && source.version.CompareMinor(target) != 0 {
		return false
	}

	return source.version.ComparePrerelease(target) != -1
}

func (source *tilde) IsPrerelease() bool {
	return source.version.IsPrerelease()
}

func (source *tilde) String() string {
	return fmt.Sprintf("~%s", source.version.String())
}
