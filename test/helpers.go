package test

import "testing"

func AssertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}

func AssertSameInt(t testing.TB, want, got int, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}
