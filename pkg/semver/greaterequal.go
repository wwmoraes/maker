package semver

import (
	"fmt"
	"strings"
)

type greaterEqual struct {
	version PartialVersion
}

// NewGreaterEqual matches versions higher or equal than itself
func NewGreaterEqual(versionStr string) (Constraint, error) {
	saneVersionStr := strings.Trim(strings.TrimPrefix(versionStr, ">="), " ")
	if len(saneVersionStr) == 0 || isXRange(saneVersionStr) {
		return nil, &ParseError{
			Func:  "NewGreaterEqual",
			Input: versionStr,
			Err:   ErrInvalidVersion,
		}
	}

	version, err := NewPartialVersion(saneVersionStr)
	if err != nil {
		return nil, &ParseError{
			Func:  "NewGreaterEqual",
			Input: versionStr,
			Err:   err,
		}
	}

	return NewGreaterEqualWith(version), nil
}

// NewGreaterEqualWith matches versions higher or equal than itself
func NewGreaterEqualWith(version PartialVersion) Constraint {
	return &greaterEqual{
		version,
	}
}

func (source *greaterEqual) Match(target Version, includePrerelease bool) bool {
	switch source.version.CompareRelease(target) {
	case 1: // release is greater than
		return !target.IsPrerelease() || includePrerelease
	case 0: // release matches
		return source.version.ComparePrerelease(target) != -1
	default: // any other cases
		return false
	}
}

func (source *greaterEqual) IsPrerelease() bool {
	return source.version.IsPrerelease()
}

func (source *greaterEqual) String() string {
	return fmt.Sprintf(">=%s", source.version.String())
}
