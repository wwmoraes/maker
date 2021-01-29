package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestInvalidNewConstraint(t *testing.T) {
	constraintStrings := []string{
		"aaa",
		"a.2.3",
		"1.a.3",
		"1.2.a",
		"1.2.3.4",
		"~aaa",
		"~1.2.3.4",
		"^aaa",
		"^1.2.3.4",
		"<aaa",
		"<1.2.3.4",
		"<=aaa",
		"<=1.2.3.4",
		"=aaa",
		"=1.2.3.4",
		">=aaa",
		">=1.2.3.4",
		">aaa",
		">1.2.3.4",
		"aaa - 2.3.4",
		"1.2.3.4 - 2.3.4",
		"1.2.3 - aaa",
		"1.2.3 - 2.3.4.5",
		"1.2.3 - 2.3.4 - 3.4.5",
		"aaa || 2.3.4",
		"1.2.3.4 || 2.3.4",
		"1.2.3 || aaa",
		"1.2.3 || 2.3.4.5",
		">1.2.3 <aaa",
		">1.2.3 <2.3.4.5",
		">aaa <2.3.4",
		">1.2.3.4 <2.3.4",
	}

	executeInvalidConstraintWith(t, constraintStrings, semver.NewConstraint)
}
