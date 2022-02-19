package gobas

import (
	"bytes"
	"testing"

	"github.com/mazzegi/gobas/testutil"
)

type rawReadTestCase struct {
	name      string
	in        string
	expectErr bool
	expectOut []rawLine
}

var rawReadTestCases = []rawReadTestCase{
	{
		name: "minimal working",
		in: `
			120 PRINT "YOU NOW HAVE ";Q;" DOLLARS."
			130 PRINT
			140 GOTO 260
		`,
		expectErr: false,
		expectOut: []rawLine{
			{1, 120, `PRINT "YOU NOW HAVE ";Q;" DOLLARS."`},
			{2, 130, `PRINT`},
			{3, 140, `GOTO 260`},
		},
	},
	{
		name: "minimal working - trimm",
		in: `
			120 PRINT "YOU NOW HAVE ";Q;" DOLLARS."
			130 PRINT
			140     GOTO 260
		`,
		expectErr: false,
		expectOut: []rawLine{
			{1, 120, `PRINT "YOU NOW HAVE ";Q;" DOLLARS."`},
			{2, 130, `PRINT`},
			{3, 140, `GOTO 260`},
		},
	},
	{
		name: "empty lines",
		in: `

			120 PRINT "YOU NOW HAVE ";Q;" DOLLARS."
			130 PRINT

			140 GOTO 260
		`,
		expectErr: false,
		expectOut: []rawLine{
			{2, 120, `PRINT "YOU NOW HAVE ";Q;" DOLLARS."`},
			{3, 130, `PRINT`},
			{5, 140, `GOTO 260`},
		},
	},
	{
		name: "fail: no line-number",
		in: `
			120 PRINT "YOU NOW HAVE ";Q;" DOLLARS."
			130PRINT
			140 GOTO 260
		`,
		expectErr: true,
		expectOut: []rawLine{},
	},
	{
		name: "fail: invalid line-number",
		in: `
			120 PRINT "YOU NOW HAVE ";Q;" DOLLARS."
			x130 PRINT
			140 GOTO 260
		`,
		expectErr: true,
		expectOut: []rawLine{},
	},
}

func TestRawRead(t *testing.T) {
	for _, test := range rawReadTestCases {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBufferString(test.in)
			rls, err := rawRead(buf)
			if err != nil {
				if !test.expectErr {
					t.Fatalf("want NO error, got %v", err)
				}
			} else if !testutil.SlicesEqual(test.expectOut, rls) {
				t.Fatalf("want %v, got %v", test.expectOut, rls)
			}
		})
	}
}
