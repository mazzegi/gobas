package expr

import (
	"math"

	"github.com/pkg/errors"
)

type Lookuper interface {
	LookupVar(name string) (interface{}, error)
	CanEvalFloat(name string) bool
}

type Op string

const (
	OpPlus  Op = "PLUS"
	OpTimes Op = "TIMES"
	OpExp   Op = "EXP"
	OpLs    Op = "LS"
	OpGt    Op = "GT"
	OpEq    Op = "EQ"
	OpLsEq  Op = "LSEQ"
	OpGtEq  Op = "GTEQ"
	OpAND   Op = "AND"
	OpOR    Op = "OR"
)

func (op Op) String() string {
	switch op {
	case OpPlus:
		return "+"
	case OpTimes:
		return "*"
	case OpExp:
		return "^"
	case OpLs:
		return "<"
	case OpGt:
		return ">"
	case OpEq:
		return "="
	case OpLsEq:
		return "<="
	case OpGtEq:
		return ">="
	case OpAND:
		return "AND"
	case OpOR:
		return "OR"
	default:
		return ""
	}
}

func (op Op) Rank() int {
	switch op {
	case OpPlus:
		return 1
	case OpTimes:
		return 2
	case OpExp:
		return 3
	case OpAND, OpOR:
		return 4
	case OpLs, OpGt, OpEq, OpLsEq, OpGtEq:
		return 5
	default:
		return 0
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
	if op.Rank() == s.Op.Rank() {
		s.Evalers = append(s.Evalers, ev)
	} else if op.Rank() > s.Op.Rank() {
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
	} else { //if op < s.Op
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
		return ConvertToFloat(v)
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

		case OpLs:
			if v < nv {
				v = 1
			} else {
				v = 0
			}
		case OpGt:
			if v > nv {
				v = 1
			} else {
				v = 0
			}
		case OpEq:
			if v == nv {
				v = 1
			} else {
				v = 0
			}
		case OpLsEq:
			if v <= nv {
				v = 1
			} else {
				v = 0
			}
		case OpGtEq:
			if v >= nv {
				v = 1
			} else {
				v = 0
			}
		case OpAND:
			if v > 0 && nv > 0 {
				v = 1
			} else {
				v = 0
			}
		case OpOR:
			if v > 0 || nv > 0 {
				v = 1
			} else {
				v = 0
			}
		}
	}
	return v, nil
}
