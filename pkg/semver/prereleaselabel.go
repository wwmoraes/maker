package semver

import (
	"regexp"
	"strings"
)

var prereleaseLabelRule = regexp.MustCompile(`^(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)$`)

// NewPrereleaseLabel returns a Label with the identifiers present on the src
// string splitted by dot, if the src is a valid prerelease label
func NewPrereleaseLabel(labelStr string) (Label, error) {
	saneLabelStr := strings.TrimPrefix(labelStr, "-")

	if len(saneLabelStr) == 0 {
		return nil, nil
	}

	if !prereleaseLabelRule.MatchString(saneLabelStr) {
		return nil, &ParseError{
			Func:  "NewPrereleaseLabel",
			Input: labelStr,
			Err:   ErrInvalidIdentifier,
		}
	}

	return newLabel(saneLabelStr), nil
}
