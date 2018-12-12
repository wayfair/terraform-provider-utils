package rand

import (
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// String
// -----------------------------------------------------------------------------

// Ensures the output of String returns a string of the correct length
func TestString_Length(t *testing.T) {
	testCases := []struct {
		inputLen    int
		expectedLen int
	}{
		{
			inputLen: 0,
		},
		{
			inputLen: 1,
		},
		{
			inputLen: 25,
		},
	}

	for _, testCase := range testCases {
		actualLen := len(String(testCase.inputLen, Lower))
		if actualLen != testCase.inputLen {
			t.Fatalf(
				"String did not return a string of the correct length. "+
					"Expected [%d], got [%d] for requested length [%d].",
				testCase.inputLen,
				actualLen,
				testCase.inputLen,
			)
		}
	}
}

// Helper function that ensures String panics. If the call to string does not
// panic, then the test is failed with a fatal.
func assertStringPanic(t *testing.T, n int, a string, f func(int, string) string) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf(
				"String did not panic for inputs n: [%d], alphabet: [%s]",
				n,
				a,
			)
		}
	}()
	f(n, a)
}

// Ensures string panics when given bad inputs
func TestString_Panic(t *testing.T) {
	testCases := []struct {
		inputLen      int
		inputAlphabet string
	}{
		{
			inputLen:      -1,
			inputAlphabet: "abc",
		},
		{
			inputLen:      -100,
			inputAlphabet: "abc",
		},
		{
			inputLen:      0,
			inputAlphabet: "",
		},
		{
			inputLen:      100,
			inputAlphabet: "",
		},
		{
			inputLen:      -1,
			inputAlphabet: "",
		},
	}
	for _, testCase := range testCases {
		assertStringPanic(
			t,
			testCase.inputLen,
			testCase.inputAlphabet,
			String,
		)
	}
}

// Ensures the output string returned does not contain unexpected characters.
// All of the characters should be part of the alphabet string.
func TestString_Characters(t *testing.T) {
	testCases := []struct {
		alphabet string
	}{
		{alphabet: Lower},
		{alphabet: Upper},
		{alphabet: Digit},
		{alphabet: Whitespace},
		{alphabet: Special},
	}
	for _, testCase := range testCases {
		outString := String(40, testCase.alphabet)
		outRunes := []rune(outString)
		for _, char := range outRunes {
			if !strings.ContainsRune(testCase.alphabet, char) {
				t.Fatalf(
					"String returned an invalid rune in output string. Rune [%s] "+
						"does not exist in alphabet [%s].",
					string(char),
					testCase.alphabet,
				)
			}
		}
	}
}

// -----------------------------------------------------------------------------
// IntArrayUnique
// -----------------------------------------------------------------------------

// Ensures the output array from IntArrayUnique is the expected length
func TestIntArrayUnique_Length(t *testing.T) {
	testCases := []struct {
		inputLen int
	}{
		{inputLen: 0},
		{inputLen: 1},
		{inputLen: 25},
	}
	for _, testCase := range testCases {
		actualLen := len(IntArrayUnique(testCase.inputLen))
		if actualLen != testCase.inputLen {
			t.Fatalf(
				"IntArrayUnique did not return an array of the correct length. "+
					"Expected [%d], got [%d] for requested length [%d].",
				testCase.inputLen,
				actualLen,
				testCase.inputLen,
			)
		}
	}
}

// Helper function that ensures IntArrayUnique panics. If the call to
// IntArrayUnique does not panic, then the test is failed with a fatal.
func assertIntArrayUniquePanic(t *testing.T, n int, f func(int) []int) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf(
				"IntArrayUnique did not panic for inputs n: [%d]",
				n,
			)
		}
	}()
	f(n)
}

// Ensures IntArrayUnique panics when given invalid parameters
func TestIntArrayUnique_Panic(t *testing.T) {
	testCases := []struct {
		inputLen int
	}{
		{inputLen: -1},
		{inputLen: -10},
	}
	for _, testCase := range testCases {
		assertIntArrayUniquePanic(
			t,
			testCase.inputLen,
			IntArrayUnique,
		)
	}
}
