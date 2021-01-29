package semver

import "regexp"

var fullRule = regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

// Version is implemented by fully semantic-version-compliant values
type Version interface {
	PartialVersion

	// Release returns the version string without any labels
	Release() string
}

// NewVersion creates a version value with the provided string, granted
// that it is fully semantic
func NewVersion(versionStr string) (Version, error) {
	// TODO replace regexp check with native, faster checks
	if !fullRule.MatchString(versionStr) {
		return nil, &ParseError{
			Func:  "NewVersion",
			Input: versionStr,
			Err:   ErrInvalidVersion,
		}
	}

	return decode(versionStr)
}
