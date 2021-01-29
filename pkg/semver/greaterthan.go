package semver

import (
	"fmt"
	"strings"
)

type greaterThan struct {
	version PartialVersion
}

// NewGreaterThan matches only versions higher than itself
func NewGreaterThan(versionStr string) (Constraint, error) {
	saneVersionStr := strings.Trim(strings.TrimPrefix(versionStr, ">"), " ")
	if len(saneVersionStr) == 0 || isXRange(saneVersionStr) {
		return nil, &ParseError{
			Func:  "NewGreaterThan",
			Input: versionStr,
			Err:   ErrInvalidVersion,
		}
	}

	version, err := NewPartialVersion(saneVersionStr)
	if err != nil {
		return nil, &ParseError{
			Func:  "NewGreaterThan",
			Input: versionStr,
			Err:   err,
		}
	}

	return NewGreaterThanWith(version), nil
}

// NewGreaterThanWith matches only versions higher than itself
func NewGreaterThanWith(version PartialVersion) Constraint {
	return &greaterThan{
		version,
	}
}

func (source *greaterThan) Match(target Version, includePrerelease bool) bool {
	switch source.version.CompareRelease(target) {
	case 1: // greater than
		return !target.IsPrerelease() || includePrerelease
	case 0: // release matches
		return source.IsPrerelease() && source.version.ComparePrerelease(target) == 1
	default: // any other cases
		return false
	}
}

func (source *greaterThan) IsPrerelease() bool {
	return source.version.IsPrerelease()
}

func (source *greaterThan) String() string {
	return fmt.Sprintf(">%s", source.version.String())
}
