package validation

import (
	"testing"
)

// -----------------------------------------------------------------------------
// DiffSuppressStringIgnoreCase
// -----------------------------------------------------------------------------

func TestDiffSuppressStringIgnoreCase(t *testing.T) {
	testCases := []struct {
		Old           string
		New           string
		ExpectedValue bool
	}{
		{
			Old:           "",
			New:           "",
			ExpectedValue: true,
		},
		{
			Old:           "",
			New:           "foo",
			ExpectedValue: false,
		},
		{
			Old:           "foo",
			New:           "",
			ExpectedValue: false,
		},
		{
			Old:           "foo",
			New:           "foo",
			ExpectedValue: true,
		},
		{
			Old:           "FOO",
			New:           "FOO",
			ExpectedValue: true,
		},
		{
			Old:           "Foo",
			New:           "foo",
			ExpectedValue: true,
		},
		{
			Old:           "foo",
			New:           "bar",
			ExpectedValue: false,
		},
	}

	for _, testCase := range testCases {
		actualValue := DiffSuppressStringIgnoreCase(
			// empty string for key - should not matter
			"",
			testCase.Old,
			testCase.New,
			// nil for resource data reference - should not matter
			nil,
		)
		if actualValue != testCase.ExpectedValue {
			t.Fatalf(
				"DiffSuppressStringIgnoreCase did not return the correct value. "+
					"Expected [%t], got [%t] when given [%s], [%s] as old value and "+
					"new value.",
				testCase.ExpectedValue,
				actualValue,
				testCase.Old,
				testCase.New,
			)
		}
	}
}
