package nfa

import (
	"strconv"
	"testing"
)

type test struct {
	input string
	match bool
}

func check(p NFA, tests []test, t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			m := p.Match(tt.input)
			if m != tt.match {
				t.Errorf("matching %s: got %v, want %v", strconv.Quote(tt.input), m, tt.match)
			}
		})
	}
}

func TestNFA(t *testing.T) {
	p := Seq(Star(Or(S("foo"), S("bar"))), E())

	tests := []test{
		{"foo", true},
		{"bar", true},
		{"foobar", true},
		{"farboo", false},
		{"boofar", false},
		{"barfoo", true},
		{"foofoobarfooX", false},
		{"foofoobarfoo", true},
	}

	check(p, tests, t)
}
