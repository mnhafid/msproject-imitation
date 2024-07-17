package main

import (
	"errors"
)

func splitString(s string, orig string) (string, int64, int64, float64, error) {
	var (
		v, f  int64       // integers before, after decimal point
		scale float64 = 1 // value = v + f/scale
	)
	var err error
	// The next character must be [0-9.]
	if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
		return s, v, f, scale, errors.New("time: invalid duration " + quote(orig))
	}
	// Consume [0-9]*
	pl := len(s)
	v, s, err = leadingInt(s)
	if err != nil {
		return s, v, f, scale, errors.New("time: invalid duration " + quote(orig))
	}
	pre := pl != len(s) // whether we consumed anything before a period

	// Consume (\.[0-9]*)?
	post := false
	if s != "" && s[0] == '.' {
		s = s[1:]
		pl := len(s)
		f, scale, s = leadingFraction(s)
		post = pl != len(s)
	}
	if !pre && !post {
		// no digits (e.g. ".s" or "-.s")
		return s, v, f, scale, errors.New("time: invalid duration " + quote(orig))
	}
	return s, v, f, scale, nil
}

func consumeUnit(s string, original string) (int, error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c == '.' || '0' <= c && c <= '9' {
			break
		}
	}
	if i == 0 {
		return 0, errors.New("time: missing unit in duration " + quote(original))
	}

	return i, nil
}

func splitLag(s string) (string, int64) {
	orig := s
	s, v, _, _, _ := splitString(s, orig)
	return s, v
}
