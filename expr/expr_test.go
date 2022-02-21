package expr

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

func dumpStack(s *Stack) string {
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}

func floatsEqual(f1, f2 float64) bool {
	return math.Abs(f1-f2) < 1e-6
}

func setupFuncs() *Funcs {
	fs := NewFuncs()
	fs.AddFloatFunc("sqrt", func(vs []interface{}) (float64, error) {
		var a float64
		if err := ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Sqrt(a), nil
	})
	fs.AddFloatFunc("min", func(vs []interface{}) (float64, error) {
		var a1, a2 float64
		if err := ScanArgs(vs, &a1, &a2); err != nil {
			return 0, err
		}
		return math.Min(a1, a2), nil
	})
	return fs
}

func TestFloatExpr(t *testing.T) {
	funcs := setupFuncs()
	vars := NewVars()
	vars.Add("x1", 1)
	vars.Add("x2", 2)
	vars.Add("x44", 44)

	tests := []struct {
		in        string
		failParse bool
		failEval  bool
		expect    float64
		exclusive bool
	}{
		{
			in:        "min(  (2+ 1) *4, 5.67)",
			failParse: false,
			expect:    5.67,
		},
		{
			in:        "(2+1)*4",
			failParse: false,
			expect:    12,
		},
		{
			in:        "(1+1)*4+2*3",
			failParse: false,
			expect:    14,
			//exclusive: true,
		},
		{
			in:        "(sqrt(4)+1+1)*4+2*3",
			failParse: false,
			expect:    22,
		},
		{
			in:        "1+2*3",
			failParse: false,
			expect:    7,
		},
		{
			in:        "1+2*3/3^(1+(2+5))",
			failParse: false,
			expect:    1.000914,
		},
		{
			in:        "sqrt(4)-1+2*3",
			failParse: false,
			expect:    7,
		},
		{
			in:        "((sqrt(4)+1)+1)*4+2*3",
			failParse: false,
			expect:    22,
		},
		{
			in:        "x1 + x2 + x44",
			failParse: false,
			expect:    47,
		},
	}

	skipNonExclusive := false

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			if skipNonExclusive && !test.exclusive {
				t.Log("Skip")
				return
			}

			s, err := NewParser(test.in, funcs).Parse()

			if err != nil {
				if !test.failParse {
					t.Fatalf("expect NO parse error, but got %v", err)
				}
			} else {
				v, err := s.Eval(vars, funcs)
				if err != nil {
					if !test.failEval {
						t.Fatalf("expect NO eval error, but got %v", err)
					}
				} else {
					f, err := convertToFloat(v)
					if err != nil {
						t.Fatalf("expect NO convert-float error, but got %v", err)
					}
					if !floatsEqual(f, test.expect) {
						t.Fatalf("expect %f, got %f\n%s", test.expect, f, dumpStack(s))
						//t.Fatalf("expect %f, got %f", test.expect, f)
					}
				}
			}
		})
	}
}
