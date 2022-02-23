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
	OpNotEq Op = "NEQ"
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
	case OpNotEq:
		return "<>"
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
	case OpLs, OpGt, OpEq, OpNotEq, OpLsEq, OpGtEq:
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

	// evalFloat := func(e Evaler, lu Lookuper, funcs *Funcs) (float64, error) {
	// 	v, err := e.Eval(lu, funcs)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	return ConvertToFloat(v)
	// }

	// if there is more than 1 elt on the stack it must be convertible to float64 (ops +/*/^)
	var v any
	for i, e := range s.Evalers {
		nv, err := e.Eval(lu, funcs)
		//nv, err := evalFloat(e, lu, funcs)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			v = nv
			continue
		}

		switch s.Op {
		case OpPlus:
			v, err = floatOp(v, nv, func(f1, f2 float64) float64 { return f1 + f2 })
		case OpTimes:
			v, err = floatOp(v, nv, func(f1, f2 float64) float64 { return f1 * f2 })
		case OpExp:
			v, err = floatOp(v, nv, func(f1, f2 float64) float64 { return math.Pow(f1, f2) })
		case OpLs:
			v, err = less(v, nv)
		case OpGt:
			v, err = greater(v, nv)
		case OpEq:
			v, err = equal(v, nv)
		case OpNotEq:
			v, err = notEqual(v, nv)
		case OpLsEq:
			v, err = lessEqual(v, nv)
		case OpGtEq:
			v, err = greaterEqual(v, nv)
		case OpAND:
			v, err = floatOp(v, nv, func(f1, f2 float64) float64 {
				if f1 > 0 && f2 > 0 {
					return 1
				}
				return 0
			})
		case OpOR:
			v, err = floatOp(v, nv, func(f1, f2 float64) float64 {
				if f1 > 0 || f2 > 0 {
					return 1
				}
				return 0
			})
		}
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}
