package semver

import (
	"regexp"
	"strings"
)

var buildLabelRule = regexp.MustCompile(`^(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*)$`)

// NewBuildLabel returns a Label with the identifiers present on the src
// string splitted by dot, if the src is a valid build label
func NewBuildLabel(labelStr string) (Label, error) {
	saneLabelStr := strings.TrimPrefix(labelStr, "+")

	if len(saneLabelStr) == 0 {
		return nil, nil
	}

	if !buildLabelRule.MatchString(saneLabelStr) {
		return nil, &ParseError{
			Func:  "NewBuildLabel",
			Input: labelStr,
			Err:   ErrInvalidIdentifier,
		}
	}

	return newLabel(saneLabelStr), nil
}
