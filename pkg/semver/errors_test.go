package semver_test

import (
	"fmt"
	"testing"

	"github.com/wwmoraes/maker/pkg/semver"
)

func TestParseError(t *testing.T) {
	funcName := "TestParseError"
	inputStr := "lorem ipsum"
	err := semver.ErrInvalidVersion
	outputTemplate := "[semver.%s] parsing '%s': %s"
	wantStr := fmt.Sprintf(outputTemplate, funcName, inputStr, err.Error())

	parseErr := &semver.ParseError{
		Func:  funcName,
		Input: inputStr,
		Err:   err,
	}

	gotStr := parseErr.Error()

	if gotStr != wantStr {
		t.Fatalf("got [%++v], want [%++v]", gotStr, wantStr)
	}
}
