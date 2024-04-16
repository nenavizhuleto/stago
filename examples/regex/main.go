package main

import (
	"fmt"

	"github.com/nenavizhuleto/stago/fsm"
)

const (
	ASCII_MIN_DISPLAYABLE = 32
	ASCII_MAX_DISPLAYABLE = 126
)

type Pattern string

func (p Pattern) Compile() *fsm.FSM[int, rune] {
	tt := make(fsm.TT[int, rune])

	ps := string(p)

	head := 0

	for head < len(ps) {
		n := len(tt)
		c := rune(ps[head])
		switch c {
		// w/o support of '.*'
		case '*':
			// guard for '*' only
			if head-1 < 0 {
				break
			}

			// should not initialize new state at this point
			// only modify previous state.
			prev := rune(ps[head-1])
			for i := ASCII_MIN_DISPLAYABLE; i <= ASCII_MAX_DISPLAYABLE; i++ {
				if rune(i) == prev {
					// take previous state and loopback 'zero or more' character
					tt[n-1][prev] = n - 1
				} else {
					// if we encounter something else, continue matching
					tt[n-1][rune(i)] = n + 1
				}
			}
		case '.':
			tt.InitState(n)
			// 32..=126 - ASCII Displayable characters
			for i := ASCII_MIN_DISPLAYABLE; i <= ASCII_MAX_DISPLAYABLE; i++ {
				tt[n][rune(i)] = n + 1
			}
		default:
			tt.InitState(n)
			tt[n][c] = n + 1
		}

		head += 1
	}

	return fsm.New(0, tt)
}

type Regex struct {
	fsm *fsm.FSM[int, rune]
}

func NewRegex(pattern Pattern) *Regex {
	return &Regex{
		fsm: pattern.Compile(),
	}
}

func (r *Regex) Match(str string) bool {
	defer r.fsm.Reset(0)

	for _, c := range str {
		if !r.fsm.Next(c) {
			return false
		}
	}

	return true
}

func main() {
	pattern := Pattern("cba*")
	matches := []string{
		"abc",
		"aabc",
		"bc",
		"bbcd",
		"2bc",
		"bac",
		"cba",
		"cbaa",
		"cab",
	}

	regex := NewRegex(pattern)

	for _, match := range matches {
		fmt.Printf("%s -> %v\n", match, regex.Match(match))
	}
}
