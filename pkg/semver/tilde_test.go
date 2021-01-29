package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewTilde_FullVersion(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"~2.0.1",
		},
		versions: []versionScenario{
			{false, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.1"},
			{false, "1.1.0"},
			{false, "1.1.1"},
			{false, "2.0.0"},
			{false, "2.0.0-alpha"},
			{false, "2.0.0-alpha.1"},
			{false, "2.0.1-alpha"},
			{false, "2.0.1-rc.1"},
			{true, "2.0.1"},
			{true, "2.0.2"},
			{false, "2.1.0"},
			{false, "2.1.1"},
			{false, "3.0.0"},
			{false, "3.0.1"},
			{false, "3.1.0"},
			{false, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewTilde)
}

func TestValidNewTilde_AllowPatch(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"~2.0.0",
			"~2.0.x",
			"~2.0.*",
			"~2.0",
		},
		versions: []versionScenario{
			{false, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.1"},
			{false, "1.1.0"},
			{false, "1.1.1"},
			{true, "2.0.0"},
			{false, "2.0.0-alpha"},
			{false, "2.0.0-alpha.1"},
			{true, "2.0.1"},
			{true, "2.0.2"},
			{false, "2.1.0"},
			{false, "2.1.1"},
			{false, "3.0.0"},
			{false, "3.0.1"},
			{false, "3.1.0"},
			{false, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewTilde)
}

func TestValidNewTilde_AllowMinor(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"~2.x",
			"~2.X",
			"~2.*",
			"~2",
		},
		versions: []versionScenario{
			{false, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.1"},
			{false, "1.1.0"},
			{false, "1.1.1"},
			{true, "2.0.0"},
			{false, "2.0.0-alpha"},
			{false, "2.0.0-alpha.1"},
			{true, "2.0.1"},
			{true, "2.0.2"},
			{true, "2.1.0"},
			{true, "2.1.1"},
			{false, "3.0.0"},
			{false, "3.0.1"},
			{false, "3.1.0"},
			{false, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewTilde)
}

func TestValidNewTilde_Prerelease(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"~2.0.0-rc.1",
		},
		versions: []versionScenario{
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.0"},
			{false, "1.0.1"},
			{false, "1.1.0"},
			{false, "1.1.1"},
			{false, "2.0.0-alpha"},
			{false, "2.0.0-alpha.1"},
			{false, "2.0.0-rc.0"},
			{true, "2.0.0-rc.1"},
			{true, "2.0.0-rc.2"},
			{true, "2.0.0"},
			{true, "2.0.1"},
			{false, "2.1.0"},
			{false, "2.1.1"},
			{false, "3.0.0"},
			{false, "3.0.1"},
			{false, "3.1.0"},
			{false, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewTilde)
}

func TestInvalidNewTilde(t *testing.T) {
	constraints := []string{
		"~",
		"~*",
		"~x",
		"~X",
	}

	executeInvalidConstraintWith(t, constraints, semver.NewTilde)
}
