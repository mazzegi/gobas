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
	Op      Op
	Evalers []Evaler
	Encaps  bool
}

func NewStack(op Op, ev Evaler) *Stack {
	s := &Stack{
		Op:      op,
		Evalers: []Evaler{ev},
	}
	return s
}

func (s *Stack) Encapsulate() {
	s.Encaps = true
}

func (s *Stack) Push(op Op, ev Evaler) Evaler {
	if op == s.Op {
		s.Evalers = append(s.Evalers, ev)
	} else if op > s.Op {
		if s.Encaps {
			//copy this stack and append new evaler
			sub := &Stack{
				Op:      s.Op,
				Evalers: s.Evalers,
				Encaps:  s.Encaps,
			}
			s.Op = op
			s.Evalers = []Evaler{sub, ev}
			s.Encaps = false
		} else {
			last := s.Evalers[len(s.Evalers)-1]
			last = last.Push(op, ev)
			s.Evalers[len(s.Evalers)-1] = last
		}
	} else if op < s.Op {
		relts := s.Evalers
		rop := s.Op
		s.Op = op
		s.Evalers = []Evaler{
			&Stack{
				Evalers: relts,
				Op:      rop,
			},
			ev,
		}
	}
	return s
}

func (s *Stack) CanEvalFloat(lu Lookuper, funcs *Funcs) bool {
	if len(s.Evalers) == 0 {
		return false
	}
	return s.Evalers[0].CanEvalFloat(lu, funcs)
}

func (s *Stack) Eval(lu Lookuper, funcs *Funcs) (interface{}, error) {
	if len(s.Evalers) == 0 {
		return nil, errors.Errorf("no elements on the stack")
	}
	if len(s.Evalers) == 1 {
		return s.Evalers[0].Eval(lu, funcs)
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
	for i, e := range s.Evalers {
		nv, err := evalFloat(e, lu, funcs)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			v = nv
			continue
		}
		switch s.Op {
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
