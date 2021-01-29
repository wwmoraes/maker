package semver

import (
	"fmt"
	"strings"
)

type equal struct {
	version PartialVersion
}

// NewEqual returns a constraint from the version string that only matches the
// exact same version
// e.g. 1.1.1 or 1.1.1-alpha.1
func NewEqual(constraintStr string) (Constraint, error) {
	saneVersionStr := strings.Trim(strings.TrimPrefix(constraintStr, "="), " ")
	if len(saneVersionStr) == 0 {
		return nil, &ParseError{
			Func:  "NewEqual",
			Input: constraintStr,
			Err:   ErrInvalidVersion,
		}
	}

	version, err := NewPartialVersion(saneVersionStr)
	if err != nil {
		return nil, &ParseError{
			Func:  "NewEqual",
			Input: constraintStr,
			Err:   err,
		}
	}

	if version.Major() == -1 {
		return nil, &ParseError{
			Func:  "NewEqual",
			Input: constraintStr,
			Err:   ErrInvalidVersion,
		}
	}

	return NewEqualWith(version), nil
}

// NewEqualWith returns a constraint from the version value that only matches
// the exact same version
// e.g. 1.1.1 or 1.1.1-alpha.1
func NewEqualWith(version PartialVersion) Constraint {
	// TODO validate version provided

	return &equal{
		version,
	}
}

func (source *equal) Match(target Version, includePrerelease bool) bool {
	if target.IsPrerelease() && !includePrerelease {
		// do not match prereleases
		if !source.IsPrerelease() {
			return false
		}
	}

	// compare the major
	if source.version.CompareMajor(target) != 0 {
		return false
	}

	// compare the minor if it is set
	if source.version.Minor() >= 0 && source.version.CompareMinor(target) != 0 {
		return false
	}

	// compare the patch if it is set
	if source.version.Patch() >= 0 && source.version.ComparePatch(target) != 0 {
		return false
	}

	// compare the prerelease label
	return source.version.ComparePrerelease(target) == 0
}

func (source *equal) IsPrerelease() bool {
	return source.version.IsPrerelease()
}

func (source *equal) String() string {
	return fmt.Sprintf("=%s", source.version.String())
}
