package semver

import (
	"fmt"
	"strings"
)

type lessEqual struct {
	version PartialVersion
}

// NewLessEqual matches versions lower or equal than itself
func NewLessEqual(versionStr string) (Constraint, error) {
	saneVersionStr := strings.Trim(strings.TrimPrefix(versionStr, "<="), " ")
	if len(saneVersionStr) == 0 || isXRange(saneVersionStr) {
		return nil, &ParseError{
			Func:  "NewLessEqual",
			Input: versionStr,
			Err:   ErrInvalidVersion,
		}
	}

	version, err := NewPartialVersion(saneVersionStr)
	if err != nil {
		return nil, &ParseError{
			Func:  "NewLessEqual",
			Input: versionStr,
			Err:   err,
		}
	}

	return NewLessEqualWith(version), nil
}

// NewLessEqualWith matches versions lower or equal than itself
func NewLessEqualWith(version PartialVersion) Constraint {
	return &lessEqual{
		version,
	}
}

func (source *lessEqual) Match(target Version, includePrerelease bool) bool {
	switch source.version.CompareRelease(target) {
	case -1: // release is less than
		return !target.IsPrerelease() || includePrerelease
	case 0: // release matches
		// do not include prereleases
		if !source.IsPrerelease() && target.IsPrerelease() {
			return false
		}
		return source.version.ComparePrerelease(target) != 1
	default: // any other cases
		return false
	}
}

func (source *lessEqual) IsPrerelease() bool {
	return source.version.IsPrerelease()
}

func (source *lessEqual) String() string {
	return fmt.Sprintf("<=%s", source.version.String())
}
