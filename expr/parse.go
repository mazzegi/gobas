package expr

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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

	//vars := NewVars()
	var stack *Stack
	push := func(ev Evaler) error {
		defer func() {
			curr = ""
			lastOp = 0
		}()
		if stack == nil {
			if evStack, ok := ev.(*Stack); ok {
				stack = evStack
			} else {
				stack = NewStack(OpPlus, ev)
			}
			//TODO: Handle minus as first op

			// switch lastOp {
			// case minus:
			// 	// if !ev.CanEvalFloat(vars, p.funcs) {
			// 	// 	return errors.Errorf("cannot eval float")
			// 	// }
			// 	//TODO: Handle minus
			// 	stack = NewStack(OpPlus, ev)
			// 	return nil
			// case times, div, exp:
			// 	return errors.Errorf("invalid operator at beginning of expr")
			// }
			// stack = NewStack(OpPlus, ev)
			return nil
		}
		if lastOp == 0 {
			return errors.Errorf("missing operator")
		}
		// if !ev.CanEvalFloat(vars, p.funcs) {
		// 	return errors.Errorf("cannot eval float")
		// }

		switch lastOp {
		case plus:
			stack.Push(OpPlus, ev)
		case minus:
			stack.Push(OpPlus, MakeFloatFuncEvaler(ev, func(f float64) float64 { return -f }))
		case times:
			stack.Push(OpTimes, ev)
		case div:
			stack.Push(OpTimes, MakeFloatFuncEvaler(ev, func(f float64) float64 { return -1.0 / f }))
		case exp:
			stack.Push(OpExp, ev)
		}
		return nil
	}

	pushCurr := func() error {
		curr = strings.TrimSpace(curr)
		if curr == "" {
			return nil
		}
		var ev Evaler

		if strings.HasPrefix(curr, `"`) && strings.HasSuffix(curr, `"`) {
			curr = strings.TrimPrefix(curr, `"`)
			curr = strings.TrimSuffix(curr, `"`)
			ev = StringEvaler(curr)
		} else if f, err := strconv.ParseFloat(curr, 10); err == nil {
			ev = MakeNumberEvaler(f)
		} else {
			ev = VarEvaler(curr)
		}

		return push(ev)
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
				vp.Encapsulate()
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
