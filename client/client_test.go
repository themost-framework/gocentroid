package client

import (
	"regexp"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	name := "Gladys"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg := Hello("Gladys")
	if !want.MatchString(msg) {
		t.Fatalf(`Hello("Gladys") = %q, %v`, msg, want)
	}
}
