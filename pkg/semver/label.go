package semver

import (
	"strings"
)

// TODO expose interfaces for both prerelease and build instead of a struct

// Label represents a set of identifiers that form a prerelease or build version
// metadata, separated by dot
type Label []string

// newLabel returns a Label with the identifiers present on the src string
// splitted by dot
func newLabel(src string) Label {
	return strings.Split(src, ".")
}

// String returns the label identifiers as a semantic version label string
func (source Label) String() string {
	return strings.Join(source, ".")
}

// Compare returns 1 if the target is higher, -1 if lower, or 0 if equal. It
// also compares the tag version if the name matches
func (source Label) Compare(target Label) int {
	sourceLength := len(source)
	targetLength := len(target)

	// no target label vs labelled source means that target is higher, as
	// 1.0.0-alpha < 1.0.0
	if sourceLength > 0 && targetLength == 0 {
		return 1
	}

	// no source label vs labelled target means that source is higher, as
	// 1.0.0 > 1.0.0-alpha
	if targetLength > 0 && sourceLength == 0 {
		return -1
	}

	diff := 0
	// compare identifiers up to the lowest common length
	for i, j := 0, 0; i < sourceLength && j < targetLength; i, j = i+1, j+1 {
		diff = strings.Compare((target)[j], (source)[i])
		if diff != 0 {
			return diff
		}
	}

	// at this point, both labels are equal up to their lowest common length, e.g.
	// with source "alpha.1" and target "alpha.1.1", the "alpha.1" identifiers
	// were compared. To tiebreak, the amount of identifiers is checked, and the
	// label with more is considered higher (as it is more specialized)
	diff = targetLength - sourceLength
	if diff > 0 {
		return 1
	} else if diff < 0 {
		return -1
	}

	return 0
}
