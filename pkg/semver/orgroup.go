package semver

import "strings"

type orGroup struct {
	constraints []Constraint
}

// NewOrGroup matches any of the provided semantic version rules
func NewOrGroup(constraintStr string) (Constraint, error) {
	constraintParts := strings.SplitN(constraintStr, "||", 2)
	constraints := make([]Constraint, len(constraintParts))

	for index, constraintPartStr := range constraintParts {
		constraintPartStr = strings.Trim(constraintPartStr, " ")
		if len(constraintPartStr) == 0 || isXRange(constraintPartStr) {
			return nil, &ParseError{
				Func:  "NewOrGroup",
				Input: constraintStr,
				Err:   ErrInvalidVersion,
			}
		}

		constraint, err := NewConstraint(constraintPartStr)
		if err != nil {
			return nil, &ParseError{
				Func:  "NewOrGroup",
				Input: constraintStr,
				Err:   err,
			}
		}

		constraints[index] = constraint
	}

	return NewOrGroupWith(constraints...), nil
}

// NewOrGroupWith matches any of the provided semantic version rules
func NewOrGroupWith(constraints ...Constraint) Constraint {
	return &orGroup{
		constraints,
	}
}

func (source *orGroup) Match(version Version, includePrerelease bool) bool {
	includePrerelease = includePrerelease || source.IsPrerelease()

	for _, constraint := range source.constraints {
		if constraint.Match(version, includePrerelease && !constraint.IsPrerelease()) {
			return true
		}
	}

	return false
}

func (source *orGroup) IsPrerelease() bool {
	for _, constraint := range source.constraints {
		if constraint.IsPrerelease() {
			return true
		}
	}

	return false
}

func (source *orGroup) String() string {
	versionStrings := make([]string, len(source.constraints))

	for index, constraint := range source.constraints {
		versionStrings[index] = constraint.String()
	}

	return strings.Join(versionStrings, " || ")
}
