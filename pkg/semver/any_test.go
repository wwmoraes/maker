package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewAny(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"*",
			"x",
			"X",
			"",
		},
		versions: []versionScenario{
			{true, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{true, "1.0.1"},
			{true, "1.1.0"},
			{true, "1.1.1"},
			{true, "2.0.0"},
			{false, "2.0.0-alpha"},
			{false, "2.0.0-alpha.1"},
			{true, "2.0.1"},
			{true, "2.1.0"},
			{true, "2.1.1"},
			{true, "3.0.0"},
			{true, "3.0.1"},
			{true, "3.1.0"},
			{true, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewAny)
}

func TestInvalidNewAny(t *testing.T) {
	constraints := []string{
		"*.1",
		"x.1",
		"X.1",
		"_",
	}

	executeInvalidConstraintWith(t, constraints, semver.NewAny)
}

func TestNewAny_IsPrerelease(t *testing.T) {
	constraintStrings := []string{
		"*",
		"x",
		"X",
		"",
	}

	for _, constraintStr := range constraintStrings {
		t.Run(
			constraintStr,
			runnableConstraintIsPrereleaseWith(constraintStr, false, semver.NewAny),
		)
	}
}
