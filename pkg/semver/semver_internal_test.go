package semver

import (
	"testing"
)

type releaseScenario struct {
	version string
	want    string
}

func runnableDecode(versionStr string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		_, err := decode(versionStr)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	}
}

func TestInvalidDecode(t *testing.T) {
	versionStrings := []string{
		"*-alpha",
		"1-alpha",
		"1.2-alpha",
		"1.*-alpha",
		"1.2.*-alpha",
		"*.2",
		"*.2.3",
		"1.*.3",
		"1.2.3-.alpha",
		"1.2.3-alpha.",
		"1.2.3+.0096f03392876480d5ee52bada55679e",
		"1.2.3+0096f03392876480d5ee52bada55679e.",
		"a.2.3",
		"a.2.3-alpha",
		"a.2.3+0096f03392876480d5ee52bada55679e",
		"1.b.3",
		"1.b.3-alpha",
		"1.b.3+0096f03392876480d5ee52bada55679e",
		"1.2.c",
		"1.2.c-alpha",
		"1.2.c+0096f03392876480d5ee52bada55679e",
		"1.2.3.d",
		"1.2.3.d-alpha",
		"1.2.3.d+0096f03392876480d5ee52bada55679e",
		"1.2.3.4",
		"1.2.3.4-alpha",
		"1.2.3.4+0096f03392876480d5ee52bada55679e",
		"1.2.3.4+0096f03392876480d5ee52bada55679e",
	}

	for _, versionStr := range versionStrings {
		t.Run(versionStr, runnableDecode(versionStr))
	}
}

func runnablePartialVersionScenario(tt releaseScenario) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		version, err := decode(tt.version)
		if err != nil {
			t.Fatalf("expected nil error, got [%++v]", err)
		}

		got := version.Release()
		if got != tt.want {
			t.Fatalf("got [%++v], want [%++v]", got, tt.want)
		}
	}
}

func TestSemverRelease(t *testing.T) {
	testCases := []releaseScenario{
		{"x", ""},
		{"X", ""},
		{"*", ""},
		{"1.x", "1"},
		{"1.X", "1"},
		{"1.*", "1"},
		{"1.2.x", "1.2"},
		{"1.2.X", "1.2"},
		{"1.2.*", "1.2"},
	}

	for _, tt := range testCases {
		t.Run(tt.version, runnablePartialVersionScenario(tt))
	}
}
