// Package test is a helper package for testing.
//
// TODO: Need to write tests for those assertions
package test

import "testing"

// AssertSameString asserts two strings are equal
func AssertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}

// AssertSameInt asserts two ints are equal
func AssertSameInt(t testing.TB, want, got int, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}
