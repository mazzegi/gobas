package expr

import (
	"fmt"
	"math"
	"testing"

	"github.com/pkg/errors"
)

func floatsEqual(f1, f2 float64) bool {
	return math.Abs(f1-f2) < 1e-6
}

func setupFuncs() *Funcs {
	fs := NewFuncs()
	fs.AddFloatFunc("sqrt", func(vs []interface{}) (float64, error) {
		if len(vs) != 1 {
			return 0, errors.Errorf("expect 1 param, got %d", len(vs))
		}
		f, err := convertToFloat(vs[0])
		if err != nil {
			return 0, err
		}
		return math.Sqrt(f), nil
	})
	return fs
}

func TestFloatExpr(t *testing.T) {
	funcs := setupFuncs()
	vars := NewVars()

	tests := []struct {
		in        string
		failParse bool
		failEval  bool
		expect    float64
	}{
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
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
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
						t.Fatalf("expect %f, got %f", test.expect, f)
					}
				}
			}
		})
	}
}
