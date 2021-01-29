package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewGreaterThan(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			">1.0.0",
		},
		versions: []versionScenario{
			{false, "0.1.0"},
			{false, "0.1.0-alpha"},
			{false, "1.0.0"},
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

	executeConstraintScenarioWith(t, scenario, semver.NewGreaterThan)
}

func TestValidNewGreaterThan_Prerelease(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			">1.0.0-rc.1",
		},
		versions: []versionScenario{
			{false, "0.1.0"},
			{false, "0.1.0-alpha"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.0-rc.0"},
			{false, "1.0.0-rc.1"},
			{true, "1.0.0-rc.2"},
			{true, "1.0.0"},
			{true, "1.0.1"},
			{true, "1.1.0"},
			{true, "1.1.1"},
			{false, "2.0.0-alpha"},
			{false, "2.0.0-alpha.1"},
			{false, "2.0.0-rc.0"},
			{false, "2.0.0-rc.1"},
			{false, "2.0.0-rc.2"},
			{true, "2.0.0"},
			{true, "2.0.1"},
			{true, "2.1.0"},
			{true, "2.1.1"},
			{true, "3.0.0"},
			{true, "3.0.1"},
			{true, "3.1.0"},
			{true, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewGreaterThan)
}
