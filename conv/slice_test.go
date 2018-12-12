package conv

import (
	"math/rand"
	"testing"
)

// ----------------------------------------------------------------------------
// InterfaceSliceToIntSlice
// ----------------------------------------------------------------------------

// Ensures when en empty interface is supplied, or the interface value does
// not nicely assert as an int, that index will contain the value 0 in the
// int slice.
func TestInterfaceSliceToIntSlice_EmptyInterfaceValueZero(t *testing.T) {
	// NOTE(ALL): rand.Int() returns positive value
	randInt := rand.Int() % 100
	output := InterfaceSliceToIntSlice(make([]interface{}, randInt))
	// ensure each index has the value zero
	for j := 0; j < randInt; j++ {
		if output[j] != 0 {
			t.Fatalf(
				"InterfaceSliceToIntSlice did not return the correct output. "+
					"Expected [0], got [%d] for value [interface{}] at index [%d]",
				output[j],
				j,
			)
		}
	}
}

// Ensures the input slice and output slices have the same length
func TestInterfaceSliceToIntSlice_SameLength(t *testing.T) {
	// NOTE(ALL): rand.Int() returns positive value
	randInt := rand.Int() % 100
	output := len(InterfaceSliceToIntSlice(make([]interface{}, randInt)))
	if output != randInt {
		t.Fatalf(
			"InterfaceSliceToIntSlice did not return an slice with the correct "+
				"length. Expected [%d], got [%d] for value [%d].",
			randInt,
			output,
			randInt,
		)
	}
	output = len(InterfaceSliceToIntSlice(make([]interface{}, 0)))
	if output != 0 {
		t.Fatalf(
			"InterfaceSliceToIntSlice did not return an slice with the correct "+
				"length. Expected [0], got [%d] for value [0].",
			output,
		)
	}
}

// ----------------------------------------------------------------------------
// InterfaceSliceToStringSlice
// ----------------------------------------------------------------------------

// Ensures when en empty interface is supplied, or the interface value does
// not nicely assert as an int, that index will contain the value 0 in the
// int slice.
func TestInterfaceSliceToStringSlice_EmptyInterfaceValueZero(t *testing.T) {
	// NOTE(ALL): rand.Int() returns positive value
	randInt := rand.Int() % 100
	output := InterfaceSliceToStringSlice(make([]interface{}, randInt))
	// ensure each index has the value zero
	for j := 0; j < randInt; j++ {
		if output[j] != "" {
			t.Fatalf(
				"InterfaceSliceToStringSlice did not return the correct output. "+
					"Expected [\"\"], got [%s] for value [interface{}] at index [%d]",
				output[j],
				j,
			)
		}
	}
}

// Ensures the input slice and output slices have the same length
func TestInterfaceSliceToStringSlice_SameLength(t *testing.T) {
	// NOTE(ALL): rand.Int() returns positive value
	randInt := rand.Int() % 100
	output := len(InterfaceSliceToStringSlice(make([]interface{}, randInt)))
	if output != randInt {
		t.Fatalf(
			"InterfaceSliceToStringSlice did not return an slice with the correct "+
				"length. Expected [%d], got [%d] for value [%d].",
			randInt,
			output,
			randInt,
		)
	}
	output = len(InterfaceSliceToStringSlice(make([]interface{}, 0)))
	if output != 0 {
		t.Fatalf(
			"InterfaceSliceToStringSlice did not return an slice with the correct "+
				"length. Expected [0], got [%d] for value [0].",
			output,
		)
	}
}
