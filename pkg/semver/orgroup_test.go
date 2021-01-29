package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewOrGroup(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"1.x || >=3.1",
		},
		versions: []versionScenario{
			{true, "1.0.0"},
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{true, "1.0.1"},
			{true, "1.1.0"},
			{true, "1.1.1"},
			{false, "2.0.0"},
			{false, "2.0.0-alpha"},
			{false, "2.0.0-alpha.1"},
			{false, "2.0.1"},
			{false, "2.1.0"},
			{false, "2.1.1"},
			{false, "3.0.0"},
			{false, "3.0.1"},
			{true, "3.1.0"},
			{true, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewOrGroup)
}

func TestValidNewOrGroup_Prerelease(t *testing.T) {
	scenario := constraintScenario{
		constraints: []string{
			"^1.1.1-rc.1 || >=3.1",
		},
		versions: []versionScenario{
			{false, "1.0.0-alpha"},
			{false, "1.0.0-alpha.1"},
			{false, "1.0.0-rc.0"},
			{false, "1.0.0-rc.1"},
			{false, "1.0.0-rc.2"},
			{false, "1.0.0"},
			{false, "1.0.1"},
			{false, "1.1.0"},
			{false, "1.1.1-alpha"},
			{false, "1.1.1-alpha.1"},
			{true, "1.1.1-rc.1"},
			{true, "1.1.1-rc.2"},
			{true, "1.1.1"},
			{false, "2.0.0-alpha"},
			{false, "2.0.0-alpha.1"},
			{false, "2.0.0-rc.0"},
			{false, "2.0.0-rc.1"},
			{false, "2.0.0-rc.2"},
			{false, "2.0.0"},
			{false, "2.0.1"},
			{false, "2.1.0"},
			{false, "2.1.1"},
			{false, "3.0.0"},
			{false, "3.0.1"},
			{true, "3.1.0"},
			{true, "3.1.1"},
		},
	}

	executeConstraintScenarioWith(t, scenario, semver.NewOrGroup)
}

func TestInvalidNewOrGroup(t *testing.T) {
	constraints := []string{
		"= || >=3.1",
		"= ||>=3.1",
		"=|| >=3.1",
		"=||>=3.1",
		" || >=3.1",
		"|| >=3.1",
		" ||>=3.1",
		"||>=3.1",
		">=3.1 || ",
		">=3.1 ||",
		">=3.1||",
		">=3.1 ||",
		">=3.1 || =",
		">=3.1 ||=",
		">=3.1|| =",
		">=3.1||=",
	}

	executeInvalidConstraintWith(t, constraints, semver.NewOrGroup)
}
