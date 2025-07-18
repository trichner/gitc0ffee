package assert

import (
	"testing"
)

func True(t testing.TB, truth bool) {
	if !truth {
		t.Fatalf("not true")
	}
}

func EqualSlice[e comparable](t testing.TB, actual, expected []e) {
	if len(actual) != len(expected) {
		t.Fatalf("expected slice of size %d but got %d", len(expected), len(actual))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("element %d not equal: %+v != %+v", i, actual, expected)
		}
	}
}

func Len[e comparable](t testing.TB, actual []e, expected int) {
	if len(actual) != expected {
		t.Fatalf("expected size %d but got %d", expected, len(actual))
	}
}

func NoError(t testing.TB, err error) {
	if err != nil {
		t.Fatalf("expected no error, got: %+v", err)
	}
}
