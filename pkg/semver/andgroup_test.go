package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewAndGroup(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			">=3.1 <=3.3",
		},
		versions: []versionScenario{
			{false, "3.0.4"},
			{true, "3.1.0"},
			{true, "3.1.4"},
			{true, "3.2.4"},
			{true, "3.3.0"},
			{true, "3.3.1"},
			{false, "3.4.0"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewAndGroup)
}

func TestValidNewAndGroup_PrereleaseLower(t *testing.T) {
	// t.Skip("# TODO implement a way to forcefully accept prereleases")

	scenario := constraintScenario{
		constraints: []string{
			">=1.1.1-rc.1 <1.3",
		},
		versions: []versionScenario{
			{false, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.1"},
			{false, "1.1.0"},
			{false, "1.1.1-rc.0"},
			{false, "1.1.1-alpha.1"},
			{true, "1.1.1-rc.1"},
			{true, "1.1.1-rc.2"},
			{true, "1.1.1"},
			{true, "1.2.0"},
			{true, "1.2.1"},
			{false, "1.2.1-rc"},
			{true, "1.2.2"},
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

	executeConstraintScenarioWith(t, scenario, semver.NewAndGroup)
}

func TestValidNewAndGroup_PrereleaseHigher(t *testing.T) {
	// t.Skip("# TODO implement a way to forcefully accept prereleases")

	scenario := constraintScenario{
		constraints: []string{
			">=1.1.1 <1.3.0-rc.1",
		},
		versions: []versionScenario{
			{false, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.1"},
			{false, "1.1.0"},
			{false, "1.1.1-rc.0"},
			{false, "1.1.1-alpha.1"},
			{false, "1.1.1-rc.1"},
			{false, "1.1.1-rc.2"},
			{true, "1.1.1"},
			{true, "1.2.0"},
			{true, "1.2.1"},
			{false, "1.2.1-rc"},
			{true, "1.2.2"},
			{true, "1.3.0-alpha.1"},
			{true, "1.3.0-rc.0"},
			{false, "1.3.0-rc.1"},
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

	executeConstraintScenarioWith(t, scenario, semver.NewAndGroup)
}

func TestInvalidNewAndGroup(t *testing.T) {
	constraints := []string{
		"= >=3.1 <",
		">=3.1 >",
	}

	executeInvalidConstraintWith(t, constraints, semver.NewAndGroup)
}
