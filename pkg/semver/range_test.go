package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewRange(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"1.1.0 - 1.2.1",
		},
		versions: []versionScenario{
			{false, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.1"},
			{true, "1.1.0"},
			{true, "1.1.1"},
			{false, "1.1.1-beta.1"},
			{true, "1.2.0"},
			{true, "1.2.1"},
			{false, "1.2.2"},
			{false, "1.3.0"},
			{false, "1.3.1"},
			{false, "1.3.2"},
			{false, "2.0.0"},
			{false, "2.0.0-alpha"},
			{false, "2.0.0-alpha.1"},
			{false, "2.0.1"},
			{false, "2.1.0"},
			{false, "2.1.1"},
			{false, "3.0.0"},
			{false, "3.0.1"},
			{false, "3.1.0"},
			{false, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewRange)
}

func TestInvalidNewRange(t *testing.T) {
	constraints := []string{
		"* - 2.0.0",
		"1.0.0 - *",
	}

	executeInvalidConstraintWith(t, constraints, semver.NewRange)
}
