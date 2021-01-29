package semver

import (
	"fmt"
	"strings"
)

type caret struct {
	version PartialVersion
}

// NewCaret creates constraint that pins down the leftmost non-zero version number
func NewCaret(constraintStr string) (Constraint, error) {
	saneVersionStr := strings.Trim(strings.TrimPrefix(constraintStr, "^"), " ")
	if len(saneVersionStr) == 0 || isXRange(saneVersionStr) {
		return nil, &ParseError{
			Func:  "NewCaret",
			Input: constraintStr,
			Err:   ErrInvalidVersion,
		}
	}

	version, err := NewPartialVersion(saneVersionStr)
	if err != nil {
		return nil, &ParseError{
			Func:  "NewCaret",
			Input: constraintStr,
			Err:   err,
		}
	}

	return NewCaretWith(version), nil
}

// NewCaretWith creates constraint that pins down the leftmost non-zero version number
func NewCaretWith(version PartialVersion) Constraint {
	return &caret{
		version,
	}
}

func (source *caret) Match(target Version, includePrerelease bool) bool {
	if target.IsPrerelease() {
		// do not match prereleases
		if !source.IsPrerelease() {
			return false
		}
	}

	major, minor, patch := source.version.Major(), source.version.Minor(), source.version.Patch()

	// target is lower than source
	if source.version.CompareRelease(target) == -1 {
		return false
	}

	// major must match always
	if source.version.CompareMajor(target) != 0 {
		return false
	}

	// minor checks
	minorDiff := source.version.CompareMinor(target)
	// minor must match if specified when major is zero
	if major == 0 && minor >= 0 && minorDiff != 0 {
		return false
	}

	patchDiff := source.version.ComparePatch(target)
	// patch must match if specified when both major and minor are zero
	if major == 0 && minor == 0 && patch >= 0 && patchDiff != 0 {
		return false
	}

	return source.version.ComparePrerelease(target) >= 0
}

func (source *caret) IsPrerelease() bool {
	return source.version.IsPrerelease()
}

func (source *caret) String() string {
	return fmt.Sprintf("^%s", source.version.String())
}
