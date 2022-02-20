package lex

import (
	"fmt"
	"testing"
)

func SlicesEqual[T comparable](ts1, ts2 []T) bool {
	if len(ts1) != len(ts2) {
		return false
	}
	for i, t1 := range ts1 {
		if t1 != ts2[i] {
			return false
		}
	}
	return true
}

func TestSplitSkipBrackets(t *testing.T) {
	tests := []struct {
		in     string
		sep    rune
		expect []string
	}{
		{
			in:     "a,b,c,d,e",
			sep:    ',',
			expect: []string{"a", "b", "c", "d", "e"},
		},
		{
			in:     `a,b,"c,d",e`,
			sep:    ',',
			expect: []string{"a", "b", `"c,d"`, "e"},
		},
		{
			in:     "a,b(,c,d),e",
			sep:    ',',
			expect: []string{"a", "b(,c,d)", "e"},
		},
		{
			in:     `a,b(",c,d"),e`,
			sep:    ',',
			expect: []string{"a", `b(",c,d")`, "e"},
		},
		{
			in:     `a,b(",c,d),"e`,
			sep:    ',',
			expect: []string{"a", `b(",c,d),"e`},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			res := splitSkipBrackets(test.in, test.sep)
			if !SlicesEqual(test.expect, res) {
				t.Fatalf("want %v, have %v", test.expect, res)
			}
		})
	}
}
