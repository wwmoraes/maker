package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewEqual_FullVersion(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"2.0.1",
			"=2.0.1",
		},
		versions: []versionScenario{
			{false, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.1"},
			{false, "1.1.0"},
			{false, "1.1.1"},
			{false, "2.0.0"},
			{false, "2.0.1-alpha"},
			{false, "2.0.1-alpha.1"},
			{false, "2.0.1-rc.1"},
			{true, "2.0.1"},
			{false, "2.1.0"},
			{false, "2.1.1"},
			{false, "3.0.0"},
			{false, "3.0.1"},
			{false, "3.1.0"},
			{false, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewEqual)
}

func TestValidNewEqual_Prerelease(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"2.0.1-alpha.2",
			"=2.0.1-alpha.2",
		},
		versions: []versionScenario{
			{false, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.1"},
			{false, "1.1.0"},
			{false, "1.1.1"},
			{false, "2.0.0"},
			{false, "2.0.1-alpha"},
			{false, "2.0.1-alpha.1"},
			{true, "2.0.1-alpha.2"},
			{false, "2.0.1-alpha.3"},
			{false, "2.0.1-rc.1"},
			{false, "2.0.1"},
			{false, "2.1.0"},
			{false, "2.1.1"},
			{false, "3.0.0"},
			{false, "3.0.1"},
			{false, "3.1.0"},
			{false, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewEqual)
}

func TestInvalidNewEqual(t *testing.T) {
	constraints := []string{
		"=*.1",
		"=x.1",
		"=X.1",
		"=_",
		"=",
		"x",
		"X",
		"*",
		"",
	}

	executeInvalidConstraintWith(t, constraints, semver.NewEqual)
}
