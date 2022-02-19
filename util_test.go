package gobas

import (
	"testing"

	"github.com/mazzegi/gobas/testutil"
)

func TestSplitOutsideQuotes(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		sep    rune
		expect []string
	}{
		{
			name:   "regular",
			in:     "a regular testcase",
			sep:    ' ',
			expect: []string{"a", "regular", "testcase"},
		},
		{
			name:   "no sep",
			in:     "a regular testcase",
			sep:    ':',
			expect: []string{"a regular testcase"},
		},
		{
			name:   "sep at start",
			in:     ":hammer",
			sep:    ':',
			expect: []string{"", "hammer"},
		},
		{
			name:   "sep at end",
			in:     "drill:",
			sep:    ':',
			expect: []string{"drill", ""},
		},
		{
			name:   "sep at start and end",
			in:     ":hammer:drill:",
			sep:    ':',
			expect: []string{"", "hammer", "drill", ""},
		},
		{
			name:   "sep in quotes",
			in:     `a hammer in "p1:p2"`,
			sep:    ':',
			expect: []string{`a hammer in "p1:p2"`},
		},
		{
			name:   "sep in quotes and outside",
			in:     `:a hammer in:"p1:p2":`,
			sep:    ':',
			expect: []string{"", `a hammer in`, `"p1:p2"`, ""},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := splitOutsideQuotes(test.in, test.sep)
			if !testutil.SlicesEqual(test.expect, res) {
				t.Fatalf("want %v, got %v", test.expect, res)
			}
		})
	}
}
