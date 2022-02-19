package lex

import (
	"strconv"
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

func (mt MatchTarget) Eval(s string) (interface{}, error) {
	if mt.Array {
		return mt.EvalArray(s)
	}

	switch mt.Type {
	case Int:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return int(v), nil
	case String:
		return s, nil
	case Float:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, errors.Errorf("unsupported type %q", mt.Type)
	}
}

func (mt MatchTarget) EvalArray(s string) (interface{}, error) {
	var sep = ","
	if psep, ok := mt.Params["sep"]; ok {
		sep = psep
	}
	sl := strings.Split(s, sep)

	switch mt.Type {
	case Int:
		return convertStrings(sl, func(s string) (int, error) {
			v, err := strconv.ParseInt(s, 10, 64)
			return int(v), err
		})
	case String:
		return sl, nil
	case Float:
		return convertStrings(sl, func(s string) (float64, error) {
			v, err := strconv.ParseFloat(s, 64)
			return v, err
		})
	default:
		return nil, errors.Errorf("unsupported type %q", mt.Type)
	}
}

func convertStrings[T any](sl []string, parseFnc func(s string) (T, error)) ([]T, error) {
	var ts []T
	for _, s := range sl {
		t, err := parseFnc(s)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}
