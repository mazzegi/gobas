package gobas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mazzegi/gobas/expr"
	"github.com/pkg/errors"
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

type Data struct {
	pos  int
	vars []interface{}
}

func (d *Data) Add(v interface{}) {
	d.vars = append(d.vars, v)
}

func (d *Data) Read() (interface{}, error) {
	if d.pos >= len(d.vars) {
		return nil, errors.Errorf("data is empty")
	}
	v := d.vars[d.pos]
	if d.pos < len(d.vars)-1 {
		d.pos++
	}
	return v, nil
}

func (d *Data) Restore() {
	d.pos = 0
}

func (s *State) Run() {
	vars := expr.NewVars()
	funcs := BuiltinFuncs()
	arrays := map[string]any{}
	forStates := []forState{}
	data := &Data{}

	for _, line := range s.lines {
		for _, stmt := range line.stmts {
			if dataStmt, ok := stmt.(DATA); ok {
				for _, c := range dataStmt.Consts {
					c = strings.Trim(c, `"`)
					data.Add(c)
				}
			}
		}
	}

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
	var extraStmts []Stmt

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
					firstStmtIdx = fs.stmtIdx + 1
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
					s.Errorf("line %d: IF eval %q: %v", line.num, stmt.Expr.Raw, err)
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
				val, err := stmt.Expr.Stack.Eval(vars, funcs)
				if err != nil {
					s.Errorf("line %d: IF eval %q: %v", line.num, stmt.Expr.Raw, err)
					return
				}
				if s.boolVal(val) {
					extraStmts = append(extraStmts, stmt.Stmts)
				}
			case IFELSESTMT:
				val, err := stmt.Expr.Stack.Eval(vars, funcs)
				if err != nil {
					s.Errorf("line %d: eval %q", line.num, stmt.Expr.Raw)
					return
				}
				if s.boolVal(val) {
					extraStmts = append(extraStmts, stmt.Stmts)
				} else {
					extraStmts = append(extraStmts, stmt.ElseStmts)
				}
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
					s.Errorf("line %d: LET eval %q: %v", line.num, stmt.Expr.Raw, err)
					return
				}
				vars.Add(stmt.Var, val)
			case ONGOSUB:
				val, err := stmt.Expr.Stack.EvalFloat(vars, funcs)
				if err != nil {
					s.Errorf("line %d: ONGOSUB eval-float %q: %v", line.num, stmt.Expr.Raw, err)
					return
				}
				ix := int(val) - 1
				if ix < 0 || ix >= len(stmt.Lines) {
					s.Errorf("line %d: ONGOSUB invalid index %d", line.num, ix+1)
					return
				}
				ln := stmt.Lines[ix]
				nextIdx := s.findLineIdx(ln)
				if nextIdx < 0 {
					s.Errorf("line %d: ONGOSUB: no such line %d", line.num, ln)
					return
				}
				gosubFromIdx = s.currIdx
				s.currIdx = nextIdx
				continue mainloop
			case ONGOTO:
				val, err := stmt.Expr.Stack.EvalFloat(vars, funcs)
				if err != nil {
					s.Errorf("line %d: ONGOTO eval-float %q: %v", line.num, stmt.Expr.Raw, err)
					return
				}
				ix := int(math.Round(val)) - 1
				if ix < 0 || ix >= len(stmt.Lines) {
					s.Errorf("line %d: ONGOTO invalid index %d", line.num, ix+1)
					return
				}
				ln := stmt.Lines[ix]
				nextIdx := s.findLineIdx(ln)
				if nextIdx < 0 {
					s.Errorf("line %d: ONGOTO: no such line %d", line.num, ln)
					return
				}
				s.currIdx = nextIdx
				continue mainloop
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
				for _, varName := range stmt.Vars {
					dv, err := data.Read()
					if err != nil {
						s.Errorf("line %d: READ : %v", line.num, err)
						return
					}
					if !isArray(varName) {
						vars.Add(varName, dv)
						continue
					}

					ad := mustParseArray(varName)
					va, ok := arrays[ad.Var]
					if !ok {
						s.Errorf("line %d: READ no such array %q", line.num, ad.Var)
						return
					}

					var cs []int
					for _, dim := range ad.Dimensions {
						f, err := dim.Stack.EvalFloat(vars, funcs)
						if err != nil {
							s.Errorf("line %d: READ: eval-float: %v", line.num, err)
							return
						}
						cs = append(cs, int(f))
					}
					switch a := va.(type) {
					case *Array[string]:
						str, err := expr.ConvertToString(dv)
						if err != nil {
							s.Errorf("line %d: cannot convert %T to string", line.num, dv)
							return
						}
						err = a.Set(cs, str)
						if err != nil {
							s.Errorf("line %d: set array %v: %v", line.num, cs, err)
							return
						}
					case *Array[float64]:
						f, err := expr.ConvertToFloat(dv)
						if err != nil {
							s.Errorf("line %d: cannot convert %T to string", line.num, dv)
							return
						}
						err = a.Set(cs, f)
						if err != nil {
							s.Errorf("line %d: set array %v: %v", line.num, cs, err)
							return
						}
					}
				}
			case REM:
			case RESTORE:
				data.Restore()
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
					s.Errorf("line %d: eval %q: %v", line.num, stmt.Expr.Raw, err)
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
				case *Array[string]:
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
				case *Array[float64]:
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
