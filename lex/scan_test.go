package lex

import (
	"fmt"
	"testing"
)

func paramsEqual(ps1, ps2 Params) bool {
	for k1, v1 := range ps1 {
		v2, ok := ps2[k1]
		if !ok {
			return false
		}
		if v1 != v2 {
			return false
		}
	}
	for k2, v2 := range ps2 {
		v1, ok := ps1[k2]
		if !ok {
			return false
		}
		if v2 != v1 {
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

func TestScan(t *testing.T) {
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
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			ps, err := Scan(test.pattern, test.input)
			if err != nil {
				if !test.fail {
					t.Fatalf("expect NOT to fail, but got %v", err)
				}
			} else {
				if !paramsEqual(test.params, ps) {
					t.Fatalf("want %v, have %v", test.params, ps)
				}
			}
		})
	}
}
