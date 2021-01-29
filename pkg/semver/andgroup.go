package semver

import (
	"strings"
)

type andGroup struct {
	constraints []Constraint
}

// NewAndGroup matches all the provided semantic version rules
func NewAndGroup(constraintStr string) (Constraint, error) {
	constraintParts := strings.Split(constraintStr, " ")
	constraints := make([]Constraint, len(constraintParts))

	for index, constraintStr := range constraintParts {
		constraint, err := NewConstraint(constraintStr)
		if err != nil {
			return nil, &ParseError{
				Func:  "NewAndGroup",
				Input: constraintStr,
				Err:   err,
			}
		}

		constraints[index] = constraint
	}

	return NewAndGroupWith(constraints...), nil
}

// NewAndGroupWith matches all the provided semantic version rules
func NewAndGroupWith(constraints ...Constraint) Constraint {
	return &andGroup{
		constraints,
	}
}

func (source *andGroup) Match(target Version, includePrerelease bool) bool {
	includePrerelease = includePrerelease || source.IsPrerelease()

	for _, constraint := range source.constraints {
		if !constraint.Match(target, includePrerelease && !constraint.IsPrerelease()) {
			return false
		}
	}

	return true
}

func (source *andGroup) IsPrerelease() bool {
	for _, constraint := range source.constraints {
		if constraint.IsPrerelease() {
			return true
		}
	}

	return false
}

func (source *andGroup) String() string {
	versionStrings := make([]string, len(source.constraints))

	for index, constraint := range source.constraints {
		versionStrings[index] = constraint.String()
	}

	return strings.Join(versionStrings, " ")
}
