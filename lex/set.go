package lex

import (
	"strings"

	"github.com/pkg/errors"
)

func NewSet() *Set {
	return &Set{
		patterns: map[string]*Pattern{},
	}
}

type Set struct {
	patterns map[string]*Pattern
}

func (s *Set) Add(name string, pattern string) error {
	p, err := ParsePattern(pattern)
	if err != nil {
		return err
	}
	s.patterns[name] = p
	return nil
}

func (s *Set) Eval(input string) (Params, string, error) {
	for name, p := range s.patterns {
		if p.Prefix != "" && strings.HasPrefix(input, p.Prefix) {
			if ps, err := p.Eval(input); err == nil {
				return ps, name, nil
			}
		}
	}
	// none of those with prefix worked - try those without
	for name, p := range s.patterns {
		if p.Prefix != "" {
			continue
		}
		if ps, err := p.Eval(input); err == nil {
			return ps, name, nil
		}
	}
	return Params{}, "", errors.Errorf("found no pattern matching input %q", input)
}
