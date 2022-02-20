package lex

import (
	"fmt"
	"testing"
)

func TestSetEval(t *testing.T) {
	// setup set
	set := NewSet()
	set.Add("goto", "GOTO {line:int}")
	set.Add("on-gosub", "ON {expr:string} GOSUB {lines:[]int}")
	set.Add("on-goto", "ON {expr:string} GOTO {lines:[]int}")
	set.Add("def", "DEF {fnc:string}={expr:string}")
	set.Add("for-step", "FOR {ivar:string}={iexpr:string} TO {toexpr:string} STEP {stepexpr:string}")
	set.Add("for", "FOR {ivar:string}={iexpr:string} TO {toexpr:string}")
	set.Add("assign", "{var:string}={expr:string}")

	tests := []struct {
		input  string
		fail   bool
		params Params
	}{
		{
			input: "GOTO 22",
			fail:  false,
			params: Params{
				"line": 22,
			},
		},
		{
			input: "ON X GOSUB 100,200",
			fail:  false,
			params: Params{
				"expr":  "X",
				"lines": []int{100, 200},
			},
		},
		{
			input: "ON AY GOTO 21,23",
			fail:  false,
			params: Params{
				"expr":  "AY",
				"lines": []int{21, 23},
			},
		},
		{
			input: "DEF FNA(X)=SIN(x/57.3)",
			fail:  false,
			params: Params{
				"fnc":  "FNA(X)",
				"expr": "SIN(x/57.3)",
			},
		},
		{
			input: "DEF FNA(X) = SIN(x/57.3)",
			fail:  false,
			params: Params{
				"fnc":  "FNA(X)",
				"expr": "SIN(x/57.3)",
			},
		},
		{
			input: "FOR j = 2 TO N STEP 3 ",
			fail:  false,
			params: Params{
				"ivar":     "j",
				"iexpr":    "2",
				"toexpr":   "N",
				"stepexpr": "3",
			},
		},
		{
			input: "FOR ab = 45 TO 112",
			fail:  false,
			params: Params{
				"ivar":   "ab",
				"iexpr":  "45",
				"toexpr": "112",
			},
		},
		{
			input: "W21$(AB,82) = COS(i*57.3) - 1",
			fail:  false,
			params: Params{
				"var":  "W21$(AB,82)",
				"expr": "COS(i*57.3) - 1",
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			ps, _, err := set.Eval(test.input)
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
