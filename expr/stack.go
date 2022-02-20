package expr

import (
	"math"

	"github.com/pkg/errors"
)

type Lookuper interface {
	LookupVar(name string) (interface{}, error)
	CanEvalFloat(name string) bool
}

type Op int

const (
	OpPlus Op = iota
	OpTimes
	OpExp
)

func (op Op) String() string {
	switch op {
	case OpPlus:
		return "+"
	case OpTimes:
		return "*"
	default:
		return "^"
	}
}

type Evaler interface {
	Push(op Op, ev Evaler) Evaler
	Eval(lu Lookuper, funcs *Funcs) (interface{}, error)
	CanEvalFloat(lu Lookuper, funcs *Funcs) bool
}

type Stack struct {
	op      Op
	evalers []Evaler
}

func NewStack(op Op, ev Evaler) *Stack {
	s := &Stack{
		op:      op,
		evalers: []Evaler{ev},
	}
	return s
}

func (s *Stack) Push(op Op, ev Evaler) Evaler {
	if op == s.op {
		s.evalers = append(s.evalers, ev)
	} else if op > s.op {
		last := s.evalers[len(s.evalers)-1]
		last = last.Push(op, ev)
		s.evalers[len(s.evalers)-1] = last
	} else if op < s.op {
		relts := s.evalers
		rop := s.op
		s.op = op
		s.evalers = []Evaler{
			&Stack{
				evalers: relts,
				op:      rop,
			},
			ev,
		}
	}
	return s
}

func (s *Stack) CanEvalFloat(lu Lookuper, funcs *Funcs) bool {
	if len(s.evalers) == 0 {
		return false
	}
	return s.evalers[0].CanEvalFloat(lu, funcs)
}

func (s *Stack) Eval(lu Lookuper, funcs *Funcs) (interface{}, error) {
	if len(s.evalers) == 0 {
		return nil, errors.Errorf("no elements on the stack")
	}
	if len(s.evalers) == 1 {
		return s.evalers[0].Eval(lu, funcs)
	}

	evalFloat := func(e Evaler, lu Lookuper, funcs *Funcs) (float64, error) {
		v, err := e.Eval(lu, funcs)
		if err != nil {
			return 0, err
		}
		return convertToFloat(v)
	}

	// if there is more than 1 elt on the stack it must be convertible to float64 (ops +/*/^)
	var v float64
	for i, e := range s.evalers {
		nv, err := evalFloat(e, lu, funcs)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			v = nv
			continue
		}
		switch s.op {
		case OpPlus:
			v += nv
		case OpTimes:
			v *= nv
		case OpExp:
			v = math.Pow(v, nv)
		}
	}
	return v, nil
}
