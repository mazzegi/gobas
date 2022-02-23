package gobas

import (
	"fmt"
	"testing"

	"github.com/mazzegi/gobas/testutil"
)

func inputsEqual(in1, in2 INPUT) bool {
	if in1.Msg != in2.Msg {
		return false
	}
	if in1.Semicolon != in2.Semicolon {
		return false
	}
	return testutil.SlicesEqual(in1.Vars, in2.Vars)
}

func TestInput(t *testing.T) {
	tests := []struct {
		raw    string
		expect INPUT
	}{
		{
			raw: ` "please input" ; a, b`,
			expect: INPUT{
				Msg:       "please input",
				Semicolon: true,
				Vars:      []string{"a", "b"},
			},
		},
		{
			raw: ` a, b`,
			expect: INPUT{
				Msg:       "",
				Semicolon: false,
				Vars:      []string{"a", "b"},
			},
		},
		{
			raw: ` "please input"a, b`,
			expect: INPUT{
				Msg:       "please input",
				Semicolon: false,
				Vars:      []string{"a", "b"},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			res := mustParseInput(test.raw)
			if !inputsEqual(test.expect, res) {
				t.Fatalf("want %v, got %v", test.expect, res)
			}
		})
	}
}
