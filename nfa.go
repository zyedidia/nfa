package nfa

import (
	"unicode/utf8"
)

type state struct {
	final   bool
	onRune  map[rune][]*state
	onEmpty []*state
}

func newState() *state {
	return &state{
		onRune: make(map[rune][]*state),
	}
}

// Match 'str' from this state and return whether or not the match is
// successful.
func (s *state) Match(str string) bool {
	match, _ := s.match(str, []*state{})
	return match
}

func (s *state) runeEdge(r rune, next *state) {
	s.onRune[r] = append(s.onRune[r], next)
}

func (s *state) emptyEdge(next *state) {
	s.onEmpty = append(s.onEmpty, next)
}

func (s *state) match(str string, visited []*state) (bool, []*state) {
	for _, v := range visited {
		if v == s {
			return false, visited
		}
	}

	visited = append(visited, s)

	if len(str) == 0 {
		if s.final {
			return true, visited
		}

		for _, next := range s.onEmpty {
			var m bool
			m, visited = next.match("", visited)
			if m {
				return true, visited
			}
		}
		return false, visited
	}

	r, size := utf8.DecodeRuneInString(str)

	for _, next := range s.onRune[r] {
		if next.Match(str[size:]) {
			return true, visited
		}
	}

	for _, next := range s.onEmpty {
		var m bool
		m, visited = next.match(str, visited)
		if m {
			return true, visited
		}
	}

	return false, visited
}

type NFA struct {
	entry *state
	exit  *state
}

func newNFA(entry, exit *state) NFA {
	return NFA{
		entry: entry,
		exit:  exit,
	}
}

// Match returns whether this NFA accepts the given string.
func (n *NFA) Match(str string) bool {
	return n.entry.Match(str)
}

// R returns an NFA that matches the single rune 'r'.
func R(r rune) NFA {
	entry := newState()
	exit := newState()
	exit.final = true
	entry.runeEdge(r, exit)
	return newNFA(entry, exit)
}

// S returns an NFA that matches the string 's'.
func S(s string) NFA {
	if len(s) == 0 {
		return E()
	}
	r, size := utf8.DecodeRuneInString(s)
	return Seq(R(r), S(s[size:]))
}

// E returns an empty NFA.
func E() NFA {
	entry := newState()
	exit := newState()
	exit.final = true
	entry.emptyEdge(exit)
	return newNFA(entry, exit)
}

// Star returns an NFA that matches 'nfa' zero or more times.
func Star(nfa NFA) NFA {
	nfa.exit.emptyEdge(nfa.entry)
	nfa.entry.emptyEdge(nfa.exit)
	return nfa
}

func seq(first, second NFA) NFA {
	first.exit.final = false
	second.exit.final = true
	first.exit.emptyEdge(second.entry)
	return newNFA(first.entry, second.exit)
}

// Seq returns an NFA that matches all given nfas concatenated in a sequence.
func Seq(rexps ...NFA) NFA {
	if len(rexps) == 0 {
		return E()
	}

	exp := rexps[0]
	for i := 1; i < len(rexps); i++ {
		exp = seq(exp, rexps[i])
	}
	return exp
}

func or(choice1, choice2 NFA) NFA {
	choice1.exit.final = false
	choice2.exit.final = false
	entry := newState()
	exit := newState()
	exit.final = true
	entry.emptyEdge(choice1.entry)
	entry.emptyEdge(choice2.entry)
	choice1.exit.emptyEdge(exit)
	choice2.exit.emptyEdge(exit)
	return newNFA(entry, exit)
}

// Or returns an NFA that matches the alternation of all given nfas.
func Or(rexps ...NFA) NFA {
	if len(rexps) == 0 {
		return E()
	}

	exp := rexps[0]
	for i := 1; i < len(rexps); i++ {
		exp = or(exp, rexps[i])
	}
	return exp
}
