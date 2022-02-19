package lex

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	Int    = "int"
	String = "string"
	Float  = "float"
	Bool   = "bool"
)

type MatchTarget struct {
	Name   string
	Array  bool
	Type   string
	Params map[string]string
}

// ParseMatchTarget parses target strings like {expr:string} , {lines:[]int}, {name:[]string?sep=;&max=5}
func ParseMatchTarget(s string) (MatchTarget, error) {
	name, rest, ok := strings.Cut(s, ":")
	if !ok {
		return MatchTarget{}, errors.Errorf("invalid match-target syntax: must be <name:type?params>")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return MatchTarget{}, errors.Errorf("match-target name may not be empty")
	}
	mt := MatchTarget{
		Name:   name,
		Params: map[string]string{},
	}

	typ, query, _ := strings.Cut(strings.TrimSpace(rest), "?")
	if strings.HasPrefix(typ, "[]") {
		mt.Array = true
		typ = strings.TrimPrefix(typ, "[]")
	}
	if typ == "" {
		return MatchTarget{}, errors.Errorf("match-target type may not be empty")
	}
	mt.Type = typ

	query = strings.TrimSpace(query)
	if query != "" {
		qsl := strings.Split(query, "&")
		for _, qs := range qsl {
			key, val, ok := strings.Cut(qs, "=")
			if !ok {
				return MatchTarget{}, errors.Errorf("invalid query parameter %q", qs)
			}
			key = strings.TrimSpace(key)
			val = strings.TrimSpace(val)
			if key == "" {
				return MatchTarget{}, errors.Errorf("invalid query parameter %q (key is empty)", qs)
			}
			if val == "" {
				return MatchTarget{}, errors.Errorf("invalid query parameter %q (val is empty)", qs)
			}
			mt.Params[key] = val
		}
	}
	return mt, nil
}
