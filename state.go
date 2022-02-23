package gobas

import (
	"fmt"
	"strconv"

	"github.com/mazzegi/gobas/expr"
)

type Line struct {
	num   int
	stmts []Stmt
}

type State struct {
	currIdx int
	lines   []Line
}

func (s *State) findLineIdx(num int) int {
	for i, l := range s.lines {
		if l.num == num {
			return i
		}
	}
	return -1
}

func (s *State) Out(v interface{}) {
	fmt.Print(v)
}

func (s *State) Outln(v interface{}) {
	fmt.Println(v)
}

func (s *State) Errorf(pattern string, args ...interface{}) {
	fmt.Printf("ERROR: "+pattern+"\n", args...)
}

func (s *State) Outfln(pattern string, args ...interface{}) {
	fmt.Printf(pattern+"\n", args...)
}

func (s *State) Outf(pattern string, args ...interface{}) {
	fmt.Printf(pattern, args...)
}

func (s *State) boolVal(v interface{}) bool {
	f, _ := expr.ConvertToFloat(v)
	return f > 0
}

type forState struct {
	lineIdx int
	stmtIdx int
	varName string
	toValue float64
	step    float64
}

func (s *State) Run() {
	vars := expr.NewVars()
	funcs := BuiltinFuncs()
	arrays := map[string]any{}
	forStates := []forState{}

	findForState := func(varName string) (forState, bool) {
		if len(forStates) == 0 {
			return forState{}, false
		}
		if varName == "" {
			return forStates[len(forStates)-1], true
		}

		for i := len(forStates) - 1; i >= 0; i-- {
			if forStates[i].varName == varName {
				return forStates[i], true
			}
		}
		return forState{}, false
	}

	var gosubFromIdx int = -1
	var firstStmtIdx int = 0
mainloop:
	for {
		if s.currIdx >= len(s.lines) {
			s.Outln("beyond last line")
			return
		}

		line := s.lines[s.currIdx]
		for stmtIdx := firstStmtIdx; stmtIdx < len(line.stmts); stmtIdx++ {
			stmt := line.stmts[stmtIdx]
			switch stmt := stmt.(type) {
			case DATA:
				//TODO
			case DEF:
				//TODO
			case DIM:
				for _, ad := range stmt.Arrays {
					var dims []int
					for _, dim := range ad.Dimensions {
						f, err := dim.Stack.EvalFloat(vars, funcs)
						if err != nil {
							s.Errorf("line %d: DIM: eval-float: %v", line.num, err)
							return
						}
						dims = append(dims, int(f))
					}

					if IsString(ad.Var) {
						a := NewArray[string](dims)
						funcs.AddFunc(ad.Var, func(vs []interface{}) (interface{}, error) {
							cs := make([]int, len(vs))
							args := make([]interface{}, len(vs))
							for i := 0; i < len(vs); i++ {
								args[i] = &(cs[i])
							}
							if err := expr.ScanArgs(vs, args...); err != nil {
								return 0, err
							}
							return a.Get(cs)
						})
						arrays[ad.Var] = a
					} else {
						a := NewArray[float64](dims)
						funcs.AddFloatFunc(ad.Var, func(vs []interface{}) (float64, error) {
							cs := make([]int, len(vs))
							args := make([]interface{}, len(vs))
							for i := 0; i < len(vs); i++ {
								args[i] = &(cs[i])
							}
							if err := expr.ScanArgs(vs, args...); err != nil {
								return 0, err
							}
							return a.Get(cs)
						})
						arrays[ad.Var] = a
					}
				}
			case END:
				s.Outln("END")
				return
			case FOR:
				iv, err := stmt.Initial.Stack.EvalFloat(vars, funcs)
				if err != nil {
					s.Errorf("line %d: FOR: eval initial: %v", line.num, err)
					return
				}
				to, err := stmt.To.Stack.EvalFloat(vars, funcs)
				if err != nil {
					s.Errorf("line %d: FOR: eval to: %v", line.num, err)
					return
				}
				step, err := stmt.Step.Stack.EvalFloat(vars, funcs)
				if err != nil {
					s.Errorf("line %d: FOR: eval step: %v", line.num, err)
					return
				}

				vars.Add(stmt.Var, iv)
				forStates = append(forStates, forState{
					lineIdx: s.currIdx,
					stmtIdx: stmtIdx,
					varName: stmt.Var,
					toValue: to,
					step:    step,
				})

			case NEXT:
				fs, ok := findForState(stmt.Var)
				if !ok {
					s.Errorf("line %d: NEXT: found no corresponding for-state", line.num)
					return
				}
				v, err := vars.LookupVar(fs.varName)
				if err != nil {
					s.Errorf("line %d: NEXT: found no var %q", line.num, fs.varName)
				}
				f, err := expr.ConvertToFloat(v)
				if err != nil {
					s.Errorf("line %d: NEXT: eval var value %q: %v", line.num, fs.varName, err)
				}
				f++
				if f <= fs.toValue {
					vars.Add(fs.varName, f)
					s.currIdx = fs.lineIdx
					firstStmtIdx = fs.stmtIdx
					continue mainloop
				}
				//just go on

			case GOSUB:
				nextIdx := s.findLineIdx(stmt.Line)
				if nextIdx < 0 {
					s.Errorf("line %d: GOSUB: no such line %d", line.num, stmt.Line)
					return
				}
				gosubFromIdx = s.currIdx
				s.currIdx = nextIdx
				continue mainloop
			case GOTO:
				nextIdx := s.findLineIdx(stmt.Line)
				if nextIdx < 0 {
					s.Errorf("line %d: GOTO: no such line %d", line.num, stmt.Line)
					return
				}
				s.currIdx = nextIdx
				continue mainloop
			case IFLN:
				val, err := stmt.Expr.Stack.Eval(vars, funcs)
				if err != nil {
					s.Errorf("line %d: eval %q", line.num, stmt.Expr.Raw)
					return
				}
				if s.boolVal(val) {
					nextIdx := s.findLineIdx(stmt.Line)
					if nextIdx < 0 {
						s.Errorf("line %d: IFLN: no such line %d", line.num, stmt.Line)
						return
					}
					s.currIdx = nextIdx
					continue mainloop
				}
			case IFELSELN:
				val, err := stmt.Expr.Stack.Eval(vars, funcs)
				if err != nil {
					s.Errorf("line %d: eval %q", line.num, stmt.Expr.Raw)
					return
				}
				var nextIdx int
				if s.boolVal(val) {
					nextIdx = s.findLineIdx(stmt.Line)
				} else {
					nextIdx = s.findLineIdx(stmt.ElseLine)
				}
				if nextIdx < 0 {
					s.Errorf("line %d: IFLN: no such line %d", line.num, stmt.Line)
					return
				}
				s.currIdx = nextIdx
				continue mainloop
			case IFSTMT:
			case IFELSESTMT:
			case INPUT:
				s.Outf("%s? ", stmt.Msg)
			inputouter:
				for {
					var in string
					fmt.Scanln(&in)
					sl := splitOutsideQuotes(in, ',')
					if len(sl) != len(stmt.Vars) {
						s.Outfln("invalid input count %d: need %d", len(sl), len(stmt.Vars))
						continue
					}
					for i, vn := range stmt.Vars {
						if IsString(vn) {
							vars.Add(vn, sl[i])
						} else {
							f, err := strconv.ParseFloat(sl[i], 64)
							if err != nil {
								s.Outfln("cannot parse %q as float", sl[i])
								continue inputouter
							}
							vars.Add(vn, f)
						}
					}
					break
				}

			case LET:
				val, err := stmt.Expr.Stack.Eval(vars, funcs)
				if err != nil {
					s.Errorf("line %d: eval %q", line.num, stmt.Expr.Raw)
					return
				}
				vars.Add(stmt.Var, val)
			case ONGOSUB:
			case ONGOTO:
			case PRINT:
				lastSemicolon := false
				for _, pi := range stmt.Items {
					lastSemicolon = false
					switch pi := pi.(type) {
					case Expr:
						val, err := pi.Stack.Eval(vars, funcs)
						if err != nil {
							s.Errorf("line %d: eval %q", line.num, pi.Raw)
							return
						}
						s.Out(val)
					case printComma:
						s.Out("\t")
					case printSemicolon:
						lastSemicolon = true
					}
				}
				if !lastSemicolon {
					s.Out("\n")
				}
			case READ:
			case REM:
			case RESTORE:
			case RETURN:
				if gosubFromIdx < 0 {
					s.Errorf("line %d: RETURN without GOSUB", line.num)
					return
				}
				s.currIdx = gosubFromIdx + 1
				gosubFromIdx = -1
				continue mainloop
			case STOP:
				s.Outln("STOP")
				return
			case ASSIGN:
				val, err := stmt.Expr.Stack.Eval(vars, funcs)
				if err != nil {
					s.Errorf("line %d: eval %q", line.num, stmt.Expr.Raw)
					return
				}
				vars.Add(stmt.Var, val)
			case ASSIGN_ARRAY:
				va, ok := arrays[stmt.Array.Var]
				if !ok {
					s.Errorf("line %d: no such array %q", line.num, stmt.Array.Var)
					return
				}
				val, err := stmt.Expr.Stack.Eval(vars, funcs)
				if err != nil {
					s.Errorf("line %d: eval %q", line.num, stmt.Expr.Raw)
					return
				}
				var cs []int
				for _, dim := range stmt.Array.Dimensions {
					f, err := dim.Stack.EvalFloat(vars, funcs)
					if err != nil {
						s.Errorf("line %d: DIM: eval-float: %v", line.num, err)
						return
					}
					cs = append(cs, int(f))
				}
				switch a := va.(type) {
				case Array[string]:
					str, err := expr.ConvertToString(val)
					if err != nil {
						s.Errorf("line %d: cannot convert %T to string", line.num, val)
						return
					}
					err = a.Set(cs, str)
					if err != nil {
						s.Errorf("line %d: set array %v: %v", line.num, cs, err)
						return
					}
				case Array[float64]:
					f, err := expr.ConvertToFloat(val)
					if err != nil {
						s.Errorf("line %d: cannot convert %T to string", line.num, val)
						return
					}
					err = a.Set(cs, f)
					if err != nil {
						s.Errorf("line %d: set array %v: %v", line.num, cs, err)
						return
					}
				}
			}
		}
		s.currIdx++
		firstStmtIdx = 0
	}
}
