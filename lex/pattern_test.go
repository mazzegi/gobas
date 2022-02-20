package lex

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func paramsEqual(ps1, ps2 Params) bool {
	for k1, v1 := range ps1 {
		v2, ok := ps2[k1]
		if !ok {
			return false
		}
		if !reflect.DeepEqual(v1, v2) {
			return false
		}
	}
	for k2, v2 := range ps2 {
		v1, ok := ps1[k2]
		if !ok {
			return false
		}
		if !reflect.DeepEqual(v2, v1) {
			return false
		}
	}
	return true
}

func TestParamsEqual(t *testing.T) {
	tests := []struct {
		ps1 Params
		ps2 Params
		eq  bool
	}{
		{
			ps1: Params{
				"line": 22,
			},
			ps2: Params{
				"line": 22,
			},
			eq: true,
		},
		{
			ps1: Params{
				"line": 22,
				"foo":  "bar",
			},
			ps2: Params{
				"line": 22,
				"foo":  "bar",
			},
			eq: true,
		},
		{
			ps1: Params{
				"line": 22,
				"foo":  "bar",
			},
			ps2: Params{
				"line": 22,
				"foo":  21,
			},
			eq: false,
		},
		{
			ps1: Params{
				"line": 22,
			},
			ps2: Params{
				"line": 22,
				"foo":  "bar",
			},
			eq: false,
		},
		{
			ps1: Params{
				"line": 22,
				"foo":  "bar",
			},
			ps2: Params{
				"line": 22,
			},
			eq: false,
		},
		{
			ps1: Params{
				"line": []int{2, 5, 7, 9},
			},
			ps2: Params{
				"line": []int{2, 5, 7, 9},
			},
			eq: true,
		},
		{
			ps1: Params{
				"line": []int{2, 5, 9},
			},
			ps2: Params{
				"line": []int{2, 5, 7, 9},
			},
			eq: false,
		},
		{
			ps1: Params{
				"line": []int{2, 5, 7, 9},
			},
			ps2: Params{
				"line": []uint8{2, 5, 7, 9},
			},
			eq: false,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			res := paramsEqual(test.ps1, test.ps2)
			if test.eq != res {
				t.Fatalf("want %t, have %t", test.eq, res)
			}
		})
	}
}

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
		// case Whitespace:
		// 	_, ok := p2.Targets[i].(Whitespace)
		// 	if !ok {
		// 		return false
		// 	}
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
					//// Whitespace{},
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
					// Whitespace{},
					MatchTarget{
						Name:   "expr",
						Array:  false,
						Type:   "string",
						Params: map[string]string{},
					},
					// Whitespace{},
					"GOSUB",
					// Whitespace{},
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
					// Whitespace{},
					MatchTarget{
						Name:   "expr",
						Array:  false,
						Type:   "string",
						Params: map[string]string{},
					},
					"GOSUB",
					// Whitespace{},
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
		{
			in:     "{name:string}{expr:string}",
			fail:   true,
			expect: &Pattern{},
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

func dumpParams(ps Params) string {
	bs, _ := json.Marshal(ps)
	return string(bs)
}

func TestPatternEval(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		fail    bool
		params  Params
	}{
		{
			pattern: "GOTO {line:int}",
			input:   "GOTO 22",
			fail:    false,
			params: Params{
				"line": 22,
			},
		},
		{
			pattern: "ON {expr:string} GOSUB {lines:[]int}",
			input:   "ON X GOSUB 100,200",
			fail:    false,
			params: Params{
				"expr":  "X",
				"lines": []int{100, 200},
			},
		},
		{
			pattern: "DEF {fnc:string}={expr:string}",
			input:   "DEF FNA(X)=SIN(x/57.3)",
			fail:    false,
			params: Params{
				"fnc":  "FNA(X)",
				"expr": "SIN(x/57.3)",
			},
		},
		{
			pattern: "DEF {fnc:string}={expr:string}",
			input:   "DEF FNA(X) = SIN(x/57.3)",
			fail:    false,
			params: Params{
				"fnc":  "FNA(X)",
				"expr": "SIN(x/57.3)",
			},
		},
		{
			pattern: "FOR {ivar:string}={iexpr:string} TO {toexpr:string} STEP {stepexpr:string}",
			input:   "FOR j = 2 TO N STEP 3 ",
			fail:    false,
			params: Params{
				"ivar":     "j",
				"iexpr":    "2",
				"toexpr":   "N",
				"stepexpr": "3",
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			p, err := ParsePattern(test.pattern)
			if err != nil {
				t.Fatalf("parse-pattern failed: %v", err)
			}
			ps, err := p.Eval(test.input)
			if err != nil {
				if !test.fail {
					t.Fatalf("expect NOT to fail, but got %v", err)
				}
			} else {
				if !paramsEqual(test.params, ps) {
					t.Fatalf("want %s, have %s", dumpParams(test.params), dumpParams(ps))
				}
			}
		})
	}
}
