package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewPrereleaseLabel(t *testing.T) {
	labels := []string{
		"",
		"alpha",
		"alpha.1",
		"alpha.2",
		"alpha.test",
		"alpha.test.1",
		"1",
		"2",
		"-",
		"-alpha",
		"-alpha.1",
		"-alpha.2",
		"-alpha.test",
		"-alpha.test.1",
		"-1",
		"-2",
		"--01",
	}

	executeValidPrereleaseLabel(t, labels)
}

func TestInvalidNewPrereleaseLabel(t *testing.T) {
	labels := []string{
		".alpha",
		".alpha.",
		"alpha.",
		"..alpha",
		"..alpha..",
		"alpha..",
		".alpha.1",
		".alpha.1.",
		"alpha.1.",
		"01",
		".",
		"-.alpha",
		"-.alpha.",
		"-alpha.",
		"-..alpha",
		"-..alpha..",
		"-alpha..",
		"-.alpha.1",
		"-.alpha.1.",
		"-alpha.1.",
		"-.",
		"-01",
	}

	executeInvalidLabelWith(t, labels, semver.NewPrereleaseLabel)
}
