package lex

import (
	"fmt"
	"testing"
)

func patternsEqual(p1, p2 *Pattern) bool {
	if len(p1.Targets) != len(p2.Targets) {
		return false
	}
	for i, t1 := range p1.Targets {
		switch t1 := t1.(type) {
		case string:
			t2, ok := p2.Targets[i].(string)
			if !ok || t1 != t2 {
				return false
			}
		case Whitespace:
			_, ok := p2.Targets[i].(Whitespace)
			if !ok {
				return false
			}
		case MatchTarget:
			t2, ok := p2.Targets[i].(MatchTarget)
			if !ok || !matchTargetsEqual(t1, t2) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func TestParsePattern(t *testing.T) {
	tests := []struct {
		in     string
		fail   bool
		expect *Pattern
	}{
		{
			in:   "GOTO {line:int}",
			fail: false,
			expect: &Pattern{
				Targets: []interface{}{
					"GOTO",
					Whitespace{},
					MatchTarget{
						Name:   "line",
						Array:  false,
						Type:   "int",
						Params: map[string]string{},
					},
				},
			},
		},
		{
			in:   "ON {expr:string} GOSUB {lines:[]int?sep=,}",
			fail: false,
			expect: &Pattern{
				Targets: []interface{}{
					"ON",
					Whitespace{},
					MatchTarget{
						Name:   "expr",
						Array:  false,
						Type:   "string",
						Params: map[string]string{},
					},
					Whitespace{},
					"GOSUB",
					Whitespace{},
					MatchTarget{
						Name:  "lines",
						Array: true,
						Type:  "int",
						Params: map[string]string{
							"sep": ",",
						},
					},
				},
			},
		},
		{
			in:   "ON {expr:string}GOSUB {lines:[]int?sep=,}",
			fail: false,
			expect: &Pattern{
				Targets: []interface{}{
					"ON",
					Whitespace{},
					MatchTarget{
						Name:   "expr",
						Array:  false,
						Type:   "string",
						Params: map[string]string{},
					},
					"GOSUB",
					Whitespace{},
					MatchTarget{
						Name:  "lines",
						Array: true,
						Type:  "int",
						Params: map[string]string{
							"sep": ",",
						},
					},
				},
			},
		},
		{
			in:     "GOTO {{line:int}",
			fail:   true,
			expect: &Pattern{},
		},
		{
			in:     "GOTO {line:int{}",
			fail:   true,
			expect: &Pattern{},
		},
		{
			in:     "GOTO {line:int}}",
			fail:   true,
			expect: &Pattern{},
		},
		{
			in:   "{name:string}({value:int})",
			fail: false,
			expect: &Pattern{
				Targets: []interface{}{
					MatchTarget{
						Name:   "name",
						Array:  false,
						Type:   "string",
						Params: map[string]string{},
					},
					"(",
					MatchTarget{
						Name:   "value",
						Array:  false,
						Type:   "int",
						Params: map[string]string{},
					},
					")",
				},
			},
		},
		{
			in:   "{name:string}={expr:string}",
			fail: false,
			expect: &Pattern{
				Targets: []interface{}{
					MatchTarget{
						Name:   "name",
						Array:  false,
						Type:   "string",
						Params: map[string]string{},
					},
					"=",
					MatchTarget{
						Name:   "expr",
						Array:  false,
						Type:   "string",
						Params: map[string]string{},
					},
				},
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			res, err := ParsePattern(test.in)
			if err != nil {
				if !test.fail {
					t.Fatalf("expect NOT to fail, but got %v", err)
				}
			} else {
				if !patternsEqual(test.expect, res) {
					t.Fatalf("want %v, have %v", test.expect, res)
				}
			}
		})
	}
}
