package semver

import "strings"

type any struct{}

// NewAny matches any version
func NewAny(constraintStr string) (Constraint, error) {
	version, err := NewPartialVersion(strings.Trim(constraintStr, "*xX"))
	if err != nil {
		return nil, &ParseError{
			Func:  "NewAny",
			Input: constraintStr,
			Err:   err,
		}
	}

	return NewAnyWith(version), nil
}

// NewAnyWith matches any version
func NewAnyWith(version PartialVersion) Constraint {
	return &any{}
}

func (source *any) Match(target Version, includePrerelease bool) bool {
	return !target.IsPrerelease()
}

func (source *any) IsPrerelease() bool {
	return false
}

func (source *any) String() string {
	return "*"
}
