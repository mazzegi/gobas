package expr

import "github.com/pkg/errors"

func NewParser(expr string, funcs *Funcs) *Parser {
	return &Parser{
		expression: expr,
		pos:        0,
		funcs:      funcs,
	}
}

type Parser struct {
	expression string
	pos        int
	funcs      *Funcs
}

const (
	plus   byte = '+'
	minus  byte = '-'
	times  byte = '*'
	div    byte = '/'
	exp    byte = '^'
	bopen  byte = '('
	bclose byte = ')'
)

func (p *Parser) Parse() (*Stack, error) {
	var curr string
	var lastOp byte

	// num := func(s string) (float64, error) {
	// 	if s == "" {
	// 		return 0, errors.Errorf("empty identifier")
	// 	}
	// 	if v, err := strconv.ParseFloat(s, 64); err == nil {
	// 		return v, nil
	// 	} else {
	// 		v, ok := e.lookup(s)
	// 		if !ok {
	// 			return 0, errors.Errorf("no such identifier %q", curr)
	// 		}
	// 		return v, nil
	// 	}
	// }
	vars := NewVars()
	var stack *Stack
	push := func(ev Evaler) error {
		defer func() {
			curr = ""
			lastOp = 0
		}()
		if stack == nil {
			switch lastOp {
			case minus:
				if !ev.CanEvalFloat(vars, p.funcs) {
					return errors.Errorf("cannot eval float")
				}
				stack = NewStack(OpPlus, ev)
				return nil
			case times, div, exp:
				return errors.Errorf("invalid operator at beginning of expr")
			}
			stack = NewStack(OpPlus, ev)
			return nil
		}
		if lastOp == 0 {
			return errors.Errorf("missing operator")
		}
		if !ev.CanEvalFloat(vars, p.funcs) {
			return errors.Errorf("cannot eval float")
		}
		Add an evaler, which may perform some ops on a float value like -/1/...

		f, err := convertToFloat(v)
		if err != nil {
			return err
		}
		switch lastOp {
		case plus:
			stack.Push(OpPlus, MakeValueEvaler(f))
		case minus:
			stack.Push(OpPlus, MakeValueEvaler(-f))
		case times:
			stack.Push(OpTimes, MakeValueEvaler(f))
		case div:
			stack.Push(OpTimes, MakeValueEvaler(1.0/f))
		case exp:
			stack.Push(OpExp, MakeValueEvaler(f))
		}
		return nil
	}

	pushCurr := func() error {
		if curr == "" {
			return nil
		}
		// this is either a number, string or variable
		if v, err := num(curr); err == nil {
			return push(v)
		} else {
			return err
		}
	}

	//fmt.Printf("eval: %q\n", e.expr)
	for p.pos < len(p.expression) {
		r := p.expression[p.pos]
		switch r {
		case bopen:
			ic, ok := findClosingBraceIdx(p.expression[p.pos+1:])
			if !ok {
				return nil, errors.Errorf("no closing brace found for open brace at %d", p.pos)
			}
			bexpr := p.expression[p.pos+1 : p.pos+1+ic]
			if curr != "" {
				// curr should be a function name
				sl := splitArgs(bexpr)
				args := []Evaler{}
				for _, s := range sl {
					vp, err := NewParser(s, p.funcs).Parse()
					if err != nil {
						return nil, err
					}
					args = append(args, vp)
				}
				push(FuncEvaler{
					Name: curr,
					Args: args,
				})
			} else {
				vp, err := NewParser(bexpr, p.funcs).Parse()
				if err != nil {
					return nil, err
				}
				push(vp)
			}
			p.pos = p.pos + 1 + ic + 1
		case bclose:
			return nil, errors.Errorf("unexpected closing brace at %d", p.pos)
		case plus, minus, times, div, exp:
			err := pushCurr()
			if err != nil {
				return nil, err
			}
			lastOp = r
			p.pos++
		default:
			curr += string(r)
			p.pos++
		}
	}
	err := pushCurr()
	if err != nil {
		return nil, err
	}
	if stack == nil {
		return nil, errors.Errorf("no elements on the eval-stack")
	}

	return stack, nil
}

func findClosingBraceIdx(s string) (int, bool) {
	nopen := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case bopen:
			nopen++
		case bclose:
			if nopen == 0 {
				return i, true
			}
			nopen--
		}
	}
	return -1, false
}

func splitArgs(expr string) []string {
	args := []string{}
	curr := ""
	nopen := 0
	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case bopen:
			nopen++
		case bclose:
			nopen--
		case ',':
			if nopen == 0 {
				args = append(args, curr)
				curr = ""
				continue
			}
		}
		curr += string(expr[i])
	}
	args = append(args, curr)
	return args
}
