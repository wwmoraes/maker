package semver_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

type versionConstructor = func(string) (semver.Version, error)
type constraintConstructor = func(string) (semver.Constraint, error)

type constraintScenario struct {
	constraints []string
	versions    []versionScenario
}

type versionScenario struct {
	want       bool
	versionStr string
}

func mustNewVersion(tb testing.TB, versionStr string) semver.Version {
	tb.Helper()

	version, err := semver.NewVersion(versionStr)

	if err != nil {
		tb.Fatalf("unexpected error, got %v", err)
	}

	if version == nil {
		tb.Fatal("expected valid version, got nil")
	}

	gotVersionStr := version.String()
	if gotVersionStr != versionStr {
		tb.Fatalf("expected [%++v], got [%++v]", versionStr, gotVersionStr)
	}

	return version
}

func mustNewPartialVersion(tb testing.TB, versionStr string) semver.PartialVersion {
	tb.Helper()

	version, err := semver.NewPartialVersion(versionStr)

	if err != nil {
		tb.Fatalf("unexpected error, got %v", err)
	}

	if version == nil {
		tb.Fatal("expected valid version, got nil")
	}

	gotVersionStr := version.String()
	if gotVersionStr != versionStr {
		tb.Fatalf("expected [%++v], got [%++v]", versionStr, gotVersionStr)
	}

	return version
}

func runnableNewVersion(versionStr string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		version := mustNewVersion(t, versionStr)

		if version == nil {
			t.Fatal("expected valid version, got nil")
		}
	}
}

func runnablePartialVersion(versionStr string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		version := mustNewPartialVersion(t, versionStr)

		if version == nil {
			t.Fatal("expected valid version, got nil")
		}
	}
}

func runnableInvalidVersion(versionStr string, constructor versionConstructor) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		version, err := semver.NewVersion(versionStr)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var parseErr *semver.ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("expected an error wrapped with ParseError, instead got a plain %++v", err.Error())
		}

		if !errors.Is(err, semver.ErrInvalidVersion) {
			t.Fatalf("expected error [%++v], got [%++v]", semver.ErrInvalidVersion, err)
		}

		if version != nil {
			t.Fatal("expected nil, got", version)
		}
	}
}

func mustNewSpecificConstraint(tb testing.TB, constraintStr string, constructor constraintConstructor) semver.Constraint {
	tb.Helper()

	constraint, err := constructor(constraintStr)
	if err != nil {
		tb.Fatal(err)
	}

	return constraint
}

func runnableConstraintMatchVersion(constraint semver.Constraint, versionStr string, want bool) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Parallel()

		target, err := semver.NewVersion(versionStr)
		if err != nil {
			t.Fatalf("unable to parse target version %v - %++v", versionStr, err.Error())
		}
		got := constraint.Match(target, false)
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func executeConstraintScenarioWith(t *testing.T, scenario constraintScenario, constructor constraintConstructor) {
	t.Helper()

	for _, constraintStr := range scenario.constraints {
		constraint := mustNewSpecificConstraint(t, constraintStr, constructor)

		for _, tt := range scenario.versions {
			t.Run(
				fmt.Sprintf("%s § %s", tt.versionStr, constraintStr),
				runnableConstraintMatchVersion(constraint, tt.versionStr, tt.want),
			)
		}
	}
}

func runnableInvalidConstraintWith(constraintStr string, constructor constraintConstructor) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Parallel()

		constraint, err := constructor(constraintStr)

		if err == nil {
			t.Fatalf("expected an error, got nil")
		}

		var parseErr *semver.ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("expected an error wrapped with ParseError, instead got a plain %++v", err.Error())
		}

		if !errors.Is(err, semver.ErrInvalidVersion) {
			t.Fatalf("expected %++v, got %++v", semver.ErrInvalidVersion, err.Error())
		}

		if constraint != nil {
			t.Fatal("expected a nil constraint")
		}
	}
}

func executeInvalidConstraintWith(t *testing.T, constraints []string, constructor constraintConstructor) {
	t.Helper()

	for _, constraintStr := range constraints {
		t.Run(
			constraintStr,
			runnableInvalidConstraintWith(constraintStr, constructor),
		)
	}
}

func runnableVersionIsPrerelease(versionStr string, want bool) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		version := mustNewVersion(t, versionStr)

		got := version.IsPrerelease()
		if got != want {
			t.Fatalf("version prerelease is %++v, expected %++v", got, want)
		}
	}
}

func runnableConstraintIsPrereleaseWith(constraintStr string, want bool, constructor constraintConstructor) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		constraint := mustNewSpecificConstraint(t, constraintStr, constructor)

		got := constraint.IsPrerelease()
		if got != want {
			t.Fatalf("version prerelease is %++v, expected %++v", got, want)
		}
	}
}
