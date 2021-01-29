package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// nullIdentifier is the initialization value for version identifiers, and is
// not semantic on purpose: if we initialized with zero (golang default), which
// is semantically valid, we'd be unable to check partial cases, such as "~1",
// which is NOT the same as "~1.0", for instance.
const nullIdentifier = -1

// isXRange returns true if the identifier is "X", "x", "*", or empty
func isXRange(identifier string) bool {
	if identifier == "" || identifier == "*" || strings.ToLower(identifier) == "x" {
		return true
	}

	return false
}

// semver represents a semantic version as per
// https://semver.org/spec/semver.0.0.html
type semver struct {
	major, minor, patch int
	prerelease, build   *Label
}

// decode parses the string into a semver value
func decode(versionStr string) (version *semver, err error) {
	// initialize with sane defaults
	version = &semver{nullIdentifier, nullIdentifier, nullIdentifier, nil, nil}

	// extract and parse the build suffix, if present
	buildParts := strings.SplitN(versionStr, "+", 2)
	remainderVersionStr := buildParts[0]
	if len(buildParts) == 2 {
		buildLabel, err := NewBuildLabel(buildParts[1])
		if err != nil {
			return nil, &ParseError{
				Func:  "decode",
				Input: versionStr,
				Err:   err,
			}
		}

		version.build = &buildLabel
	}

	// extract and parse the prerelease suffix, if present
	prereleaseParts := strings.SplitN(remainderVersionStr, "-", 2)
	remainderVersionStr = prereleaseParts[0]
	if len(prereleaseParts) == 2 {
		prereleaseLabel, err := NewPrereleaseLabel(prereleaseParts[1])
		if err != nil {
			return nil, &ParseError{
				Func:  "decode",
				Input: versionStr,
				Err:   err,
			}
		}

		version.prerelease = &prereleaseLabel
	}

	// split the version into its identifiers
	versionIdentifiers := strings.Split(remainderVersionStr, ".")
	if len(versionIdentifiers) > 3 {
		return nil, &ParseError{
			Func:  "decode",
			Input: versionStr,
			Err:   ErrInvalidVersion,
		}
	}

	// early return if major is a X range
	if isXRange(versionIdentifiers[0]) {
		// wildcards should be at the leftmost part
		if len(versionIdentifiers) > 1 {
			return nil, &ParseError{
				Func:  "decode",
				Input: versionStr,
				Err:   ErrInvalidVersion,
			}
		}
		// partial/wildcard versions must not contain labels
		if version.prerelease != nil || version.build != nil {
			return nil, &ParseError{
				Func:  "decode",
				Input: versionStr,
				Err:   ErrInvalidVersion,
			}
		}

		return version, nil
	}

	// parse major version
	version.major, err = strconv.Atoi(versionIdentifiers[0])
	if err != nil {
		return nil, &ParseError{
			Func:  "decode",
			Input: versionStr,
			Err:   err,
		}
	}

	// early return if minor is a X range
	if len(versionIdentifiers) < 2 || isXRange(versionIdentifiers[1]) {
		// wildcards should be at the rightmost part
		if len(versionIdentifiers) > 2 {
			return nil, &ParseError{
				Func:  "decode",
				Input: versionStr,
				Err:   ErrInvalidVersion,
			}
		}
		// partial/wildcard versions must not contain labels
		if version.prerelease != nil || version.build != nil {
			return nil, &ParseError{
				Func:  "decode",
				Input: versionStr,
				Err:   ErrInvalidVersion,
			}
		}

		return version, nil
	}

	// parse minor version
	version.minor, err = strconv.Atoi(versionIdentifiers[1])
	if err != nil {
		return nil, &ParseError{
			Func:  "decode",
			Input: versionStr,
			Err:   err,
		}
	}

	// early return if patch is a X range
	if len(versionIdentifiers) < 3 || isXRange(versionIdentifiers[2]) {
		// partial/wildcard versions must not contain labels
		if version.prerelease != nil || version.build != nil {
			return nil, &ParseError{
				Func:  "decode",
				Input: versionStr,
				Err:   ErrInvalidVersion,
			}
		}

		return version, nil
	}

	// parse patch version
	version.patch, err = strconv.Atoi(versionIdentifiers[2])
	if err != nil {
		return nil, &ParseError{
			Func:  "decode",
			Input: versionStr,
			Err:   err,
		}
	}

	return version, nil
}

// String returns the full version, including any labels (prerelease/build)
func (source *semver) String() string {
	var strBuilder strings.Builder

	fmt.Fprint(&strBuilder, source.Release())

	if source.IsPrerelease() {
		fmt.Fprintf(&strBuilder, "-%s", source.Prerelease().String())
	}

	if source.IsBuild() {
		fmt.Fprintf(&strBuilder, "+%s", source.Build())
	}

	return strBuilder.String()
}

// Major returns the major identifier
func (source *semver) Major() int {
	return source.major
}

// Minor returns the minor identifier
func (source *semver) Minor() int {
	return source.minor
}

// Patch returns the patch identifier
func (source *semver) Patch() int {
	return source.patch
}

// Release returns the version string without any labels
// (e.g. prerelease or build)
func (source *semver) Release() string {
	if source.major == nullIdentifier {
		return ""
	}

	var strBuilder strings.Builder

	fmt.Fprintf(&strBuilder, "%d", source.major)

	if source.minor == nullIdentifier {
		return strBuilder.String()
	}

	fmt.Fprintf(&strBuilder, ".%d", source.minor)

	if source.patch == nullIdentifier {
		return strBuilder.String()
	}

	fmt.Fprintf(&strBuilder, ".%d", source.patch)

	return strBuilder.String()
}

// Prerelease returns the prerelease label
func (source *semver) Prerelease() Label {
	if source.prerelease == nil {
		return []string{}
	}

	return *source.prerelease
}

// IsPrerelease returns true if this is a prerelease version
func (source *semver) IsPrerelease() bool {
	if source.prerelease == nil {
		return false
	}

	return len(*source.prerelease) > 0
}

// Prerelease returns the build label
func (source *semver) Build() Label {
	if source.build == nil {
		return []string{}
	}

	return *source.build
}

// IsBuild returns true if this is a build version
func (source *semver) IsBuild() bool {
	if source.build == nil {
		return false
	}

	return len(*source.build) > 0
}

// Compare returns an integer comparing two versions lexicographically.
// The result will be 0 if target == source, -1 if target < source, and +1 if
// target > source.
func (source *semver) Compare(target PartialVersion) int {
	if diff := source.CompareRelease(target); diff != 0 {
		return diff
	}

	return source.ComparePrerelease(target)
}

func (source *semver) CompareRelease(target PartialVersion) int {
	var diff int

	diff = source.CompareMajor(target)
	if diff != 0 || source.minor == -1 {
		return diff
	}

	diff = source.CompareMinor(target)
	if diff != 0 || source.patch == -1 {
		return diff
	}

	return source.ComparePatch(target)
}

// CompareMajor returns an integer comparing two major versions lexicographically.
// The result will be 0 if target == source, -1 if target < source, and +1 if
// target > source.
func (source *semver) CompareMajor(target PartialVersion) int {
	diff := target.Major() - source.Major()
	if diff > 0 {
		return 1
	} else if diff < 0 {
		return -1
	}

	return 0
}

// CompareMinor returns an integer comparing two minor versions lexicographically.
// The result will be 0 if target == source, -1 if target < source, and +1 if
// target > source.
func (source *semver) CompareMinor(target PartialVersion) int {
	diff := target.Minor() - source.Minor()
	if diff > 0 {
		return 1
	} else if diff < 0 {
		return -1
	}

	return 0
}

// ComparePatch returns an integer comparing two patch versions lexicographically.
// The result will be 0 if target == source, -1 if target < source, and +1 if
// target > source.
func (source *semver) ComparePatch(target PartialVersion) int {
	diff := target.Patch() - source.Patch()
	if diff > 0 {
		return 1
	} else if diff < 0 {
		return -1
	}

	return 0
}

// ComparePrerelease returns an integer comparing two prerelease labels lexicographically.
// The result will be 0 if target == source, -1 if target < source, and +1 if
// target > source.
func (source *semver) ComparePrerelease(target PartialVersion) int {
	sourceIsPrerelease := source.IsPrerelease()
	targetIsPrerelease := target.IsPrerelease()

	// prerelease source version and release target version
	// ∴ target is higher version (prerelease-wise)
	if sourceIsPrerelease && !targetIsPrerelease {
		return 1
	}

	// release source version and release target version
	// ∴ both versions are equal (prerelease-wise)
	if !sourceIsPrerelease && !targetIsPrerelease {
		return 0
	}

	// release source version and prerelease target version
	// ∴ target is a lower version (prerelease-wise)
	if !sourceIsPrerelease && targetIsPrerelease {
		return -1
	}

	// both source and target versions are labelled as prereleases
	// ∴ lexicographically compare labels
	targetPrerelease := target.Prerelease()
	return source.prerelease.Compare(targetPrerelease)
}
