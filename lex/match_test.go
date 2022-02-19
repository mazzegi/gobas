package lex

import (
	"fmt"
	"testing"
)

func stringMapsEqual(ps1, ps2 map[string]string) bool {
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

func matchTargetsEqual(mt1, mt2 MatchTarget) bool {
	if mt1.Name != mt2.Name ||
		mt1.Array != mt2.Array ||
		mt1.Type != mt2.Type {
		return false
	}
	return stringMapsEqual(mt1.Params, mt2.Params)
}

func TestParseMatchTarget(t *testing.T) {
	tests := []struct {
		in     string
		fail   bool
		expect MatchTarget
	}{
		{
			in:   "line:int",
			fail: false,
			expect: MatchTarget{
				Name:   "line",
				Array:  false,
				Type:   "int",
				Params: map[string]string{},
			},
		},
		{
			in:   "names:[]string?sep=;",
			fail: false,
			expect: MatchTarget{
				Name:  "names",
				Array: true,
				Type:  "string",
				Params: map[string]string{
					"sep": ";",
				},
			},
		},
		{
			in:   "names:[]string?sep=;&max=5",
			fail: false,
			expect: MatchTarget{
				Name:  "names",
				Array: true,
				Type:  "string",
				Params: map[string]string{
					"sep": ";",
					"max": "5",
				},
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			mt, err := ParseMatchTarget(test.in)
			if err != nil {
				if !test.fail {
					t.Fatalf("expect NOT to fail, but got %v", err)
				}
			} else {
				if !matchTargetsEqual(test.expect, mt) {
					t.Fatalf("want %v, have %v", test.expect, mt)
				}
			}
		})
	}
}
