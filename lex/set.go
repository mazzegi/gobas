package lex

import (
	"strings"

	"github.com/pkg/errors"
)

func NewSet() *Set {
	return &Set{
		patterns: []*Pattern{},
	}
}

type Set struct {
	patterns []*Pattern
}

func (s *Set) MustAdd(name string, pattern string) {
	err := s.Add(name, pattern)
	if err != nil {
		panic(err)
	}
}

func (s *Set) Add(name string, pattern string) error {
	p, err := ParsePatternWithName(name, pattern)
	if err != nil {
		return err
	}
	s.patterns = append(s.patterns, p)
	return nil
}

func (s *Set) Eval(input string) (Params, string, error) {
	for _, p := range s.patterns {
		if p.Prefix != "" && strings.HasPrefix(input, p.Prefix) {
			if ps, err := p.Eval(input); err == nil {
				return ps, p.Name, nil
			}
		}
	}
	// none of those with prefix worked - try those without
	for _, p := range s.patterns {
		if p.Prefix != "" {
			continue
		}
		if ps, err := p.Eval(input); err == nil {
			return ps, p.Name, nil
		}
	}
	return Params{}, "", errors.Errorf("found no pattern matching input %q", input)
}
