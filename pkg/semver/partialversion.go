package semver

import (
	"fmt"
	"regexp"
)

var partialRule = regexp.MustCompile(`^(?:[<=>~^]*)?(?:(?P<major>0|[xX\*]|[1-9]\d*)(?:\.(?P<minor>0|[xX\*]|[1-9]\d*)(?:\.(?P<patch>0|[xX\*]|[1-9]\d*)(?:-(?P<prerelease>(?:0|[xX\*]|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[xX\*]|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?)?)?)?$`)

// PartialVersion is implemented by semantic-version-compatible values that
// are incomplete or contain special syntax to match versions
type PartialVersion interface {
	fmt.Stringer

	// Major returns the major identifier
	Major() int
	// Minor returns the minor identifier
	Minor() int
	// Patch returns the patch identifier
	Patch() int

	// Prerelease returns the prerelease label
	Prerelease() Label
	// IsPrerelease returns true if this is a prerelease version
	IsPrerelease() bool

	// Prerelease returns the build label
	Build() Label
	// IsBuild returns true if this is a build version
	IsBuild() bool

	// Compare returns an integer comparing two versions lexicographically.
	// The result will be 0 if target == source, -1 if target < source, and +1 if
	// target > source.
	Compare(target PartialVersion) int
	// CompareRelease returns an integer comparing two release versions
	// lexicographically.
	// The result will be 0 if target == source, -1 if target < source, and +1 if
	// target > source.
	CompareRelease(target PartialVersion) int
	// CompareMajor returns an integer comparing two major versions
	// lexicographically.
	// The result will be 0 if target == source, -1 if target < source, and +1 if
	// target > source.
	CompareMajor(target PartialVersion) int
	// CompareMinor returns an integer comparing two minor versions
	// lexicographically.
	// The result will be 0 if target == source, -1 if target < source, and +1 if
	// target > source.
	CompareMinor(target PartialVersion) int
	// ComparePatch returns an integer comparing two patch versions
	// lexicographically.
	// The result will be 0 if target == source, -1 if target < source, and +1 if
	// target > source.
	ComparePatch(target PartialVersion) int
	// ComparePrerelease returns an integer comparing two prerelease labels lexicographically.
	// The result will be 0 if target == source, -1 if target < source, and +1 if
	// target > source.
	ComparePrerelease(target PartialVersion) int
}

// NewPartialVersion creates a partial version value with the provided string,
// granted that the fragment is semantic and supported by constraints
func NewPartialVersion(versionStr string) (PartialVersion, error) {
	if !partialRule.MatchString(versionStr) {
		return nil, &ParseError{
			Func:  "NewPartialVersion",
			Input: versionStr,
			Err:   ErrInvalidVersion,
		}
	}

	version, err := decode(versionStr)
	if err != nil {
		return nil, &ParseError{
			Func:  "NewPartialVersion",
			Input: versionStr,
			Err:   err,
		}
	}

	return version, nil
}
