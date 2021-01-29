package semver_test

import (
	"testing"
)

func TestValidNewPartialVersion(t *testing.T) {
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
			runnablePartialVersion(versionStr),
		)
	}
}
