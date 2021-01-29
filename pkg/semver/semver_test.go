package semver_test

import (
	"bytes"
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

type fullVersionScenario struct {
	releaseStr       string
	prereleaseStr    string
	buildStr         string
	wantIsPrerelease bool
	wantIsBuild      bool
}

func runnableFullVersionScenario(tt fullVersionScenario) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		completeVersionBuf := bytes.NewBufferString(tt.releaseStr)
		if len(tt.prereleaseStr) > 0 {
			completeVersionBuf.WriteString("-" + tt.prereleaseStr)
		}
		if len(tt.buildStr) > 0 {
			completeVersionBuf.WriteString("+" + tt.buildStr)
		}
		completeVersionStr := completeVersionBuf.String()

		version := mustNewVersion(t, completeVersionStr)

		gotCompleteStr := version.String()
		if gotCompleteStr != completeVersionStr {
			t.Fatalf("expected [%++v], got [%++v]", completeVersionStr, gotCompleteStr)
		}

		gotReleaseStr := version.Release()
		if gotReleaseStr != tt.releaseStr {
			t.Fatalf("expected [%++v], got [%++v]", tt.releaseStr, gotReleaseStr)
		}

		gotIsPrerelease := version.IsPrerelease()
		if gotIsPrerelease != tt.wantIsPrerelease {
			t.Fatalf("expected [%++v], got [%++v]", tt.wantIsPrerelease, gotIsPrerelease)
		}

		wantPrerelease := mustNewLabelWith(t, tt.prereleaseStr, semver.NewPrereleaseLabel)
		gotPrerelease := version.Prerelease()
		if gotPrerelease.String() != wantPrerelease.String() {
			t.Fatalf("expected [%++v], got [%++v]", wantPrerelease, gotPrerelease)
		}

		gotIsBuild := version.IsBuild()
		if gotIsBuild != tt.wantIsBuild {
			t.Fatalf("expected [%++v], got [%++v]", tt.wantIsBuild, gotIsBuild)
		}

		wantBuild := mustNewLabelWith(t, tt.buildStr, semver.NewBuildLabel)
		gotBuild := version.Build()
		if gotBuild.String() != wantBuild.String() {
			t.Fatalf("expected [%++v], got [%++v]", wantBuild, gotBuild)
		}
	}
}

func TestFullVersion(t *testing.T) {
	testCases := []fullVersionScenario{
		{
			releaseStr:       "1.2.3",
			prereleaseStr:    "",
			buildStr:         "",
			wantIsPrerelease: false,
			wantIsBuild:      false,
		},
		{
			releaseStr:       "1.2.3",
			prereleaseStr:    "alpha.1",
			buildStr:         "",
			wantIsPrerelease: true,
			wantIsBuild:      false,
		},
		{
			releaseStr:       "1.2.3",
			prereleaseStr:    "",
			buildStr:         "0096f03392876480d5ee52bada55679e",
			wantIsPrerelease: false,
			wantIsBuild:      true,
		},
		{
			releaseStr:       "1.2.3",
			prereleaseStr:    "alpha.1",
			buildStr:         "0096f03392876480d5ee52bada55679e",
			wantIsPrerelease: true,
			wantIsBuild:      true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.releaseStr, runnableFullVersionScenario(tt))
	}
}
