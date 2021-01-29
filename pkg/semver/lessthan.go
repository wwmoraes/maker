package semver

import (
	"fmt"
	"strings"
)

type lessThan struct {
	version PartialVersion
}

// NewLessThan matches only versions lower than itself
func NewLessThan(versionStr string) (Constraint, error) {
	saneVersionStr := strings.Trim(strings.TrimPrefix(versionStr, "<"), " ")
	if len(saneVersionStr) == 0 || isXRange(saneVersionStr) {
		return nil, &ParseError{
			Func:  "NewLessThan",
			Input: versionStr,
			Err:   ErrInvalidVersion,
		}
	}

	version, err := NewPartialVersion(saneVersionStr)
	if err != nil {
		return nil, &ParseError{
			Func:  "NewLessThan",
			Input: versionStr,
			Err:   err,
		}
	}

	return NewLessThanWith(version), nil
}

// NewLessThanWith matches only versions lower than itself
func NewLessThanWith(version PartialVersion) Constraint {
	return &lessThan{
		version,
	}
}

func (source *lessThan) Match(target Version, includePrerelease bool) bool {
	switch source.version.CompareRelease(target) {
	case -1: // release is less than
		return !target.IsPrerelease() || includePrerelease
	case 0: // release matches...
		// ... so the prerelease must match as well
		return source.IsPrerelease() && source.version.ComparePrerelease(target) == -1
	default: // any other cases
		return false
	}
}

func (source *lessThan) IsPrerelease() bool {
	return source.version.IsPrerelease()
}

func (source *lessThan) String() string {
	return fmt.Sprintf("<%s", source.version.String())
}
