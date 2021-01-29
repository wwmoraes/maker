package semver_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

type labelConstructor = func(string) (semver.Label, error)

func mustNewLabelWith(tb testing.TB, labelStr string, constructor labelConstructor) semver.Label {
	tb.Helper()

	label, err := constructor(labelStr)
	if err != nil {
		tb.Fatal(err)
	}

	return label
}

func executeInvalidLabelWith(t *testing.T, labels []string, constructor labelConstructor) {
	t.Helper()

	for _, labelStr := range labels {
		t.Run(
			labelStr,
			runnableInvalidLabelWith(labelStr, constructor),
		)
	}
}

func runnableInvalidLabelWith(labelStr string, constructor labelConstructor) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		label, err := constructor(labelStr)

		if label != nil {
			t.Errorf("got [%++v] label value, expected nil", label)
		}

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var parseErr *semver.ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("expected an error wrapped with ParseError, instead got a plain %++v", err.Error())
		}

		if !errors.Is(err, semver.ErrInvalidIdentifier) {
			t.Fatalf("expected [%++v], got [%++v]", semver.ErrInvalidIdentifier, err.Error())
		}
	}
}

func executeValidPrereleaseLabel(t *testing.T, labelStrings []string) {
	t.Helper()

	for _, labelStr := range labelStrings {
		t.Run(
			labelStr,
			runnableValidPrereleaseLabel(labelStr),
		)
	}
}

func runnableValidPrereleaseLabel(labelStr string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		label := mustNewLabelWith(t, labelStr, semver.NewPrereleaseLabel)

		wantLabelStr := labelStr
		if strings.HasPrefix(labelStr, "-") {
			wantLabelStr = strings.TrimPrefix(labelStr, "-")
		}

		gotLabelStr := label.String()
		if gotLabelStr != wantLabelStr {
			t.Fatalf("got [%++v], want [%++v]", gotLabelStr, wantLabelStr)
		}
	}
}

func executeValidBuildLabel(t *testing.T, labelStrings []string) {
	t.Helper()

	for _, labelStr := range labelStrings {
		t.Run(
			labelStr,
			runnableValidBuildLabel(labelStr),
		)
	}
}

func runnableValidBuildLabel(labelStr string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		label := mustNewLabelWith(t, labelStr, semver.NewBuildLabel)

		wantLabelStr := labelStr
		if strings.HasPrefix(labelStr, "+") {
			wantLabelStr = strings.TrimPrefix(labelStr, "+")
		}

		gotLabelStr := label.String()
		if gotLabelStr != wantLabelStr {
			t.Fatalf("got [%++v], want [%++v]", gotLabelStr, wantLabelStr)
		}
	}
}

func TestLabelCompare(t *testing.T) {
	testCases := []struct {
		source, target semver.Label
		want           int
	}{
		{semver.Label{}, semver.Label{"alpha"}, -1},
		{semver.Label{}, semver.Label{}, 0},
		{semver.Label{"alpha"}, semver.Label{}, 1},
		{semver.Label{"alpha", "1"}, semver.Label{"alpha", "0"}, -1},
		{semver.Label{"alpha", "0"}, semver.Label{"alpha", "0"}, 0},
		{semver.Label{"alpha", "0"}, semver.Label{"alpha", "1"}, 1},
		{semver.Label{"alpha", "1"}, semver.Label{"alpha"}, -1},
		{semver.Label{"alpha"}, semver.Label{"alpha"}, 0},
		{semver.Label{"alpha"}, semver.Label{"alpha", "1"}, 1},
	}

	for _, tt := range testCases {
		t.Run(
			fmt.Sprint(tt.source.String(), "§", tt.target.String()),
			func(t *testing.T) {
				got := tt.source.Compare(tt.target)
				if got != tt.want {
					t.Fatalf("got [%++v], want [%++v]", got, tt.want)
				}
			},
		)
	}
}
