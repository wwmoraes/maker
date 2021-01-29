package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

// TODO prerelease tests

func TestValidNewCaret_FromMajor(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"^1.2.3",
		},
		versions: []versionScenario{
			{false, "0.0.0"},
			{false, "0.1.2"},
			{false, "0.1.3"},
			{false, "0.1.4"},
			{false, "0.2.2"},
			{false, "0.2.3"},
			{false, "0.2.4"},
			{false, "0.3.2"},
			{false, "0.3.3"},
			{false, "0.3.4"},
			{false, "1.1.2"},
			{false, "1.1.2-beta"},
			{false, "1.1.3"},
			{false, "1.1.4"},
			{false, "1.2.2"},
			{true, "1.2.3"},
			{false, "1.2.3-rc.1"},
			{true, "1.2.4"},
			{false, "1.2.4-beta"},
			{true, "1.3.0"},
			{true, "1.3.2"},
			{false, "2.0.0"},
			{false, "2.0.2"},
			{false, "2.0.2-alpha"},
			{false, "2.0.3"},
			{false, "2.0.4"},
			{false, "2.1.2"},
			{false, "2.1.3"},
			{false, "2.1.4"},
			{false, "2.2.2"},
			{false, "2.2.3"},
			{false, "2.2.4"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewCaret)
}

func TestValidNewCaret_WithinMajor(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"^1.0.0",
			"^1.0.x",
			"^1.0.*",
			"^1.0",
			"^1.x",
			"^1.*",
			"^1",
		},
		versions: []versionScenario{
			{false, "0.0.0"},
			{false, "0.1.2"},
			{false, "0.1.3"},
			{false, "0.1.4"},
			{false, "0.2.2"},
			{false, "0.2.3"},
			{false, "0.2.4"},
			{false, "0.3.2"},
			{false, "0.3.3"},
			{false, "0.3.4"},
			{true, "1.1.2"},
			{false, "1.1.2-beta"},
			{true, "1.1.3"},
			{true, "1.1.4"},
			{true, "1.2.2"},
			{true, "1.2.3"},
			{false, "1.2.3-rc.1"},
			{true, "1.2.4"},
			{false, "1.2.4-beta"},
			{true, "1.3.0"},
			{true, "1.3.2"},
			{false, "2.0.0"},
			{false, "2.0.2"},
			{false, "2.0.2-alpha"},
			{false, "2.0.3"},
			{false, "2.0.4"},
			{false, "2.1.2"},
			{false, "2.1.3"},
			{false, "2.1.4"},
			{false, "2.2.2"},
			{false, "2.2.3"},
			{false, "2.2.4"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewCaret)
}

func TestValidNewCaret_FromMinor(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"^0.2.3",
		},
		versions: []versionScenario{
			{false, "0.0.0"},
			{false, "0.1.2"},
			{false, "0.1.3"},
			{false, "0.1.4"},
			{false, "0.2.2"},
			{true, "0.2.3"},
			{true, "0.2.4"},
			{false, "0.3.2"},
			{false, "0.3.3"},
			{false, "0.3.4"},
			{false, "1.1.2"},
			{false, "1.1.2-beta"},
			{false, "1.1.3"},
			{false, "1.1.4"},
			{false, "1.2.2"},
			{false, "1.2.3"},
			{false, "1.2.3-rc.1"},
			{false, "1.2.4"},
			{false, "1.2.4-beta"},
			{false, "1.3.0"},
			{false, "1.3.2"},
			{false, "2.0.0"},
			{false, "2.0.2"},
			{false, "2.0.2-alpha"},
			{false, "2.0.3"},
			{false, "2.0.4"},
			{false, "2.1.2"},
			{false, "2.1.3"},
			{false, "2.1.4"},
			{false, "2.2.2"},
			{false, "2.2.3"},
			{false, "2.2.4"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewCaret)
}

func TestValidNewCaret_WithinMinor(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"^0.2.0",
			"^0.2.x",
			"^0.2.*",
			"^0.2",
		},
		versions: []versionScenario{
			{false, "0.0.0"},
			{false, "0.0.1"},
			{false, "0.0.2"},
			{false, "0.0.2-alpha"},
			{false, "0.0.3"},
			{false, "0.0.3-beta"},
			{false, "0.0.4"},
			{false, "0.0.5"},
			{false, "0.1.2"},
			{false, "0.1.3"},
			{false, "0.1.4"},
			{true, "0.2.2"},
			{true, "0.2.3"},
			{true, "0.2.4"},
			{false, "0.3.2"},
			{false, "0.3.3"},
			{false, "0.3.4"},
			{false, "1.1.2"},
			{false, "1.1.2-beta"},
			{false, "1.1.3"},
			{false, "1.1.4"},
			{false, "1.2.2"},
			{false, "1.2.3"},
			{false, "1.2.3-rc.1"},
			{false, "1.2.4"},
			{false, "1.2.4-beta"},
			{false, "1.3.0"},
			{false, "1.3.2"},
			{false, "2.0.0"},
			{false, "2.0.2"},
			{false, "2.0.2-alpha"},
			{false, "2.0.3"},
			{false, "2.0.4"},
			{false, "2.1.2"},
			{false, "2.1.3"},
			{false, "2.1.4"},
			{false, "2.2.2"},
			{false, "2.2.3"},
			{false, "2.2.4"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewCaret)
}

func TestValidNewCaret_FromPatch(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"^0.0.1",
		},
		versions: []versionScenario{
			{false, "0.0.0"},
			{true, "0.0.1"},
			{false, "0.0.2"},
			{false, "0.0.3"},
			{false, "0.1.2"},
			{false, "0.1.3"},
			{false, "0.1.4"},
			{false, "0.2.2"},
			{false, "0.2.3"},
			{false, "0.2.4"},
			{false, "0.3.2"},
			{false, "0.3.3"},
			{false, "0.3.4"},
			{false, "1.1.0"},
			{false, "1.1.1"},
			{false, "1.1.2"},
			{false, "1.1.2-beta"},
			{false, "1.1.3"},
			{false, "1.1.4"},
			{false, "1.2.2"},
			{false, "1.2.3"},
			{false, "1.2.3-rc.1"},
			{false, "1.2.4"},
			{false, "1.2.4-beta"},
			{false, "1.3.0"},
			{false, "1.3.2"},
			{false, "2.0.0"},
			{false, "2.0.2"},
			{false, "2.0.2-alpha"},
			{false, "2.0.3"},
			{false, "2.0.4"},
			{false, "2.1.2"},
			{false, "2.1.3"},
			{false, "2.1.4"},
			{false, "2.2.2"},
			{false, "2.2.3"},
			{false, "2.2.4"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewCaret)
}

func TestValidNewCaret_ZeroMajor(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"^0.x",
			"^0.*",
			"^0",
		},
		versions: []versionScenario{
			{true, "0.0.0"},
			{true, "0.0.1"},
			{true, "0.0.2"},
			{false, "0.0.2-alpha"},
			{true, "0.0.3"},
			{false, "0.0.3-beta"},
			{true, "0.0.4"},
			{true, "0.0.5"},
			{true, "0.1.0"},
			{true, "0.1.2"},
			{true, "0.1.3"},
			{true, "0.1.4"},
			{true, "0.2.2"},
			{true, "0.2.3"},
			{true, "0.2.4"},
			{true, "0.3.2"},
			{true, "0.3.3"},
			{true, "0.3.4"},
			{false, "1.1.2"},
			{false, "1.1.2-beta"},
			{false, "1.1.3"},
			{false, "1.1.4"},
			{false, "1.2.2"},
			{false, "1.2.3"},
			{false, "1.2.3-rc.1"},
			{false, "1.2.4"},
			{false, "1.2.4-beta"},
			{false, "1.3.0"},
			{false, "1.3.2"},
			{false, "2.0.0"},
			{false, "2.0.2"},
			{false, "2.0.2-alpha"},
			{false, "2.0.3"},
			{false, "2.0.4"},
			{false, "2.1.2"},
			{false, "2.1.3"},
			{false, "2.1.4"},
			{false, "2.2.2"},
			{false, "2.2.3"},
			{false, "2.2.4"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewCaret)
}

func TestValidNewCaret_ZeroMinor(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"^0.0.x",
			"^0.0.*",
			"^0.0",
		},
		versions: []versionScenario{
			{true, "0.0.0"},
			{true, "0.0.1"},
			{true, "0.0.2"},
			{false, "0.0.2-alpha"},
			{true, "0.0.3"},
			{false, "0.0.3-beta"},
			{true, "0.0.4"},
			{true, "0.0.5"},
			{false, "0.1.0"},
			{false, "0.1.2"},
			{false, "0.1.3"},
			{false, "0.1.4"},
			{false, "0.2.2"},
			{false, "0.2.3"},
			{false, "0.2.4"},
			{false, "0.3.2"},
			{false, "0.3.3"},
			{false, "0.3.4"},
			{false, "1.1.2"},
			{false, "1.1.2-beta"},
			{false, "1.1.3"},
			{false, "1.1.4"},
			{false, "1.2.2"},
			{false, "1.2.3"},
			{false, "1.2.3-rc.1"},
			{false, "1.2.4"},
			{false, "1.2.4-beta"},
			{false, "1.3.0"},
			{false, "1.3.2"},
			{false, "2.0.0"},
			{false, "2.0.2"},
			{false, "2.0.2-alpha"},
			{false, "2.0.3"},
			{false, "2.0.4"},
			{false, "2.1.2"},
			{false, "2.1.3"},
			{false, "2.1.4"},
			{false, "2.2.2"},
			{false, "2.2.3"},
			{false, "2.2.4"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewCaret)
}

func TestInvalidNewCaret(t *testing.T) {
	constraints := []string{
		"^*.1",
		"^x.1",
		"^X.1",
		"^_",
		"^",
	}

	executeInvalidConstraintWith(t, constraints, semver.NewCaret)
}
