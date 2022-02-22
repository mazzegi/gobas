package expr

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// func NewParser(expr string, funcs *Funcs) *Parser {
// 	return &Parser{
// 		expression: expr,
// 		pos:        0,
// 		funcs:      funcs,
// 	}
// }

func NewParser(expr string) *Parser {
	return &Parser{
		expression: expr,
		pos:        0,
	}
}

type Parser struct {
	expression string
	pos        int
	//funcs      *Funcs
}

const (
	plus   = "+"
	minus  = "-"
	times  = "*"
	div    = "/"
	exp    = "^"
	bopen  = "("
	bclose = ")"
	ls     = "<"
	gt     = ">"
	eq     = "="
	lseq   = "<="
	gteq   = ">="
	and    = "AND "
	or     = "OR "
)

var ops = []string{
	lseq,
	gteq,
	and,
	or,

	plus,
	minus,
	times,
	div,
	exp,
	bopen,
	bclose,
	ls,
	gt,
	eq,
}

func isOneOf(s string, sl []string) bool {
	for _, cs := range sl {
		if s == cs {
			return true
		}
	}
	return false
}

func (p *Parser) Parse() (*Stack, error) {
	var curr string
	//var lastOp byte
	var lastOp string

	//vars := NewVars()
	var stack *Stack
	push := func(ev Evaler) error {
		defer func() {
			curr = ""
			lastOp = ""
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
		if lastOp == "" {
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
		case ls:
			stack.Push(OpLs, ev)
		case gt:
			stack.Push(OpGt, ev)
		case eq:
			stack.Push(OpEq, ev)
		case lseq:
			stack.Push(OpLsEq, ev)
		case gteq:
			stack.Push(OpGtEq, ev)
		case and:
			stack.Push(OpAND, ev)
		case or:
			stack.Push(OpOR, ev)
		}
		return nil
	}

	peekOp := func() (string, bool) {
		for _, op := range ops {
			if strings.HasPrefix(p.expression[p.pos:], op) {
				return op, true
			}
		}
		return "", false
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

	var inQuotes bool
	for p.pos < len(p.expression) {
		r := p.expression[p.pos]
		//switch r {
		switch {
		case r == '"':
			if !inQuotes {
				inQuotes = true
			} else {
				inQuotes = false
			}
			curr += string(r)
			p.pos++
		case !inQuotes && string(r) == bopen:
			ic, ok := findClosingBraceIdx(p.expression[p.pos+1:])
			if !ok {
				return nil, errors.Errorf("no closing brace found for open brace at %d", p.pos)
			}
			bexpr := p.expression[p.pos+1 : p.pos+1+ic]
			curr = strings.TrimSpace(curr)
			if curr != "" {
				// curr should be a function name
				sl := splitArgs(bexpr)
				args := []Evaler{}
				for _, s := range sl {
					vp, err := NewParser(s).Parse()
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
				vp, err := NewParser(bexpr).Parse()
				if err != nil {
					return nil, err
				}
				vp.Encapsulate()
				push(vp)
			}
			p.pos = p.pos + 1 + ic + 1
		case !inQuotes && string(r) == bclose:
			return nil, errors.Errorf("unexpected closing brace at %d", p.pos)
		// case isOneOf(string(r), []string{plus, minus, times, div, exp}):
		// 	err := pushCurr()
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	lastOp = string(r)
		// 	p.pos++
		case inQuotes:
			curr += string(r)
			p.pos++
		default:
			if op, ok := peekOp(); ok {
				err := pushCurr()
				if err != nil {
					return nil, err
				}
				lastOp = op
				p.pos += len(op)
			} else {
				curr += string(r)
				p.pos++
			}
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
		switch string(s[i]) {
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
		switch string(expr[i]) {
		case bopen:
			nopen++
		case bclose:
			nopen--
		case ",":
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
