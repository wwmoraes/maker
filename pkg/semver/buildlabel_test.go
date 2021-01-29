package semver_test

import (
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestValidNewBuildLabel(t *testing.T) {
	labels := []string{
		"",
		"0096f03392876480d5ee52bada55679e",
		"nightly",
		"nightly.1",
		"nightly.2",
		"nightly.0096f03392876480d5ee52bada55679e",
		"0096f03392876480d5ee52bada55679e.nightly",
		"+",
		"+0096f03392876480d5ee52bada55679e",
		"+nightly",
		"+nightly.1",
		"+nightly.2",
		"+nightly.0096f03392876480d5ee52bada55679e",
		"+0096f03392876480d5ee52bada55679e.nightly",
	}

	executeValidBuildLabel(t, labels)
}

func TestInvalidNewBuildLabel(t *testing.T) {
	labels := []string{
		".nightly",
		".nightly.",
		"nightly.",
		"..nightly",
		"..nightly..",
		"nightly..",
		".0096f03392876480d5ee52bada55679e",
		".0096f03392876480d5ee52bada55679e.",
		"0096f03392876480d5ee52bada55679e.",
		"..0096f03392876480d5ee52bada55679e",
		"..0096f03392876480d5ee52bada55679e..",
		"0096f03392876480d5ee52bada55679e..",
		".nightly.0096f03392876480d5ee52bada55679e",
		".nightly.0096f03392876480d5ee52bada55679e.",
		"nightly.0096f03392876480d5ee52bada55679e.",
		".",
		"+.nightly",
		"+.nightly.",
		"+nightly.",
		"+..nightly",
		"+..nightly..",
		"+nightly..",
		"+.0096f03392876480d5ee52bada55679e",
		"+.0096f03392876480d5ee52bada55679e.",
		"+0096f03392876480d5ee52bada55679e.",
		"+..0096f03392876480d5ee52bada55679e",
		"+..0096f03392876480d5ee52bada55679e..",
		"+0096f03392876480d5ee52bada55679e..",
		"+.nightly.0096f03392876480d5ee52bada55679e",
		"+.nightly.0096f03392876480d5ee52bada55679e.",
		"+nightly.0096f03392876480d5ee52bada55679e.",
		"+.",
	}

	executeInvalidLabelWith(t, labels, semver.NewBuildLabel)
}
