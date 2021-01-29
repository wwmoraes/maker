package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewVersion(t *testing.T) {
	versionStrings := []string{
		"1.2.3",
		"1.2.3-rc",
		"1.2.3-alpha",
		"1.2.3-alpha.beta",
		"1.2.3-alpha.1",
		"1.2.3-beta+8baef20a23d16b4204f5ffc6bdb11ad1",
		"1.2.3+8baef20a23d16b4204f5ffc6bdb11ad1",
	}

	for _, versionStr := range versionStrings {
		t.Run(
			versionStr,
			runnableNewVersion(versionStr),
		)
	}
}

func TestInvalidNewVersion(t *testing.T) {
	versionStrings := []string{
		"1",
		"1.2",
		"aaa",
		"a.2.3",
		"1.a.3",
		"1.2.a",
		"1.2.3.4",
		"1-alpha.1",
		"1-beta+8baef20a23d16b4204f5ffc6bdb11ad1",
		"1+8baef20a23d16b4204f5ffc6bdb11ad1",
		"1.2-alpha.1",
		"1.2-beta+8baef20a23d16b4204f5ffc6bdb11ad1",
		"1.2+8baef20a23d16b4204f5ffc6bdb11ad1",
	}

	for _, versionStr := range versionStrings {
		t.Run(
			versionStr,
			runnableInvalidVersion(versionStr, semver.NewVersion),
		)
	}
}

func TestValidNewVersion_Prerelease(t *testing.T) {
	versionStrings := []string{
		"1.2.3-rc",
		"1.2.3-alpha",
		"1.2.3-alpha.beta",
		"1.2.3-alpha.1",
		"1.2.3-beta+8baef20a23d16b4204f5ffc6bdb11ad1",
	}

	for _, versionStr := range versionStrings {
		t.Run(versionStr, runnableVersionIsPrerelease(versionStr, true))
	}
}
