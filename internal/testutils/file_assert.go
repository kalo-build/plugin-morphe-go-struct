package testutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
)

func AssertFileEquals(t *testing.T, actualPath string, expectedPath string, msgAndArgs ...any) bool {
	t.Helper()

	actualContents, readActualErr := os.ReadFile(actualPath)
	if readActualErr != nil {
		return assert.Fail(t, fmt.Sprintf("Actual could not be read:\n\tError: %s", readActualErr), msgAndArgs...)
	}

	expectedContents, readExpectedErr := os.ReadFile(expectedPath)
	if readExpectedErr != nil {
		return assert.Fail(t, fmt.Sprintf("Expected could not be read:\n\tError: %s", readExpectedErr), msgAndArgs...)
	}

	actualString := string(actualContents)
	expectedString := string(expectedContents)
	if actualString == expectedString {
		return true
	}

	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(actualString),
		B:        difflib.SplitLines(expectedString),
		FromFile: "Actual",
		FromDate: "",
		ToFile:   "Expected",
		ToDate:   "",
		Context:  1,
	})

	return assert.Fail(t, fmt.Sprintf("File contents differ: \n\nDiff:\n%s", diff), msgAndArgs...)
}

func AssertFileContentsEquals(t *testing.T, actualPath string, expectedContents []byte, msgAndArgs ...any) bool {
	t.Helper()

	actualContents, readActualErr := os.ReadFile(actualPath)
	if readActualErr != nil {
		return assert.Fail(t, fmt.Sprintf("Actual could not be read:\n\tError: %s", readActualErr), msgAndArgs...)
	}

	actualString := string(actualContents)
	expectedString := string(expectedContents)
	if actualString == expectedString {
		return true
	}

	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(actualString),
		B:        difflib.SplitLines(expectedString),
		FromFile: "Actual",
		FromDate: "",
		ToFile:   "Expected",
		ToDate:   "",
		Context:  1,
	})

	return assert.Fail(t, fmt.Sprintf("File contents differ: \n\nDiff:\n%s", diff), msgAndArgs...)
}
