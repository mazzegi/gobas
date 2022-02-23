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

func (s *State) Run() {
	vars := expr.NewVars()
	funcs := BuiltinFuncs()

	var gosubFromIdx int = -1
mainloop:
	for {
		if s.currIdx >= len(s.lines) {
			s.Outln("beyond last line")
			return
		}

		line := s.lines[s.currIdx]
		for _, stmt := range line.stmts {
			switch stmt := stmt.(type) {
			case DATA:
				//TODO
			case DEF:
				//TODO
			case DIM:
				//TODO
			case END:
				s.Outln("END")
				return
			case FOR:
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
			case NEXT:
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
			}
		}
		s.currIdx++
	}
}
