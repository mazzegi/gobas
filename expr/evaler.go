package expr

import "constraints"

type Number interface {
	constraints.Float | constraints.Integer
}

func MakeNumberEvaler[T Number](v T) NumberEvaler[T] {
	return NumberEvaler[T]{v: v}
}

type NumberEvaler[T Number] struct {
	v T
}

func (e NumberEvaler[T]) Eval(lu Lookuper, funcs *Funcs) (interface{}, error) {
	return T(e.v), nil
}

func (e NumberEvaler[T]) Push(op Op, ev Evaler) Evaler {
	s := NewStack(op, e)
	s.Push(op, ev)
	return s
}

func (e NumberEvaler[T]) CanEvalFloat(lu Lookuper, funcs *Funcs) bool {
	return true
}

//
type StringEvaler string

func (e StringEvaler) Eval(lu Lookuper, funcs *Funcs) (interface{}, error) {
	return string(e), nil
}

func (e StringEvaler) Push(op Op, ev Evaler) Evaler {
	s := NewStack(op, e)
	s.Push(op, ev)
	return s
}

func (e StringEvaler) CanEvalFloat(lu Lookuper, funcs *Funcs) bool {
	return false
}

//

type VarEvaler string

func (e VarEvaler) Eval(lu Lookuper, funcs *Funcs) (interface{}, error) {
	return lu.LookupVar(string(e))
}

func (e VarEvaler) Push(op Op, ev Evaler) Evaler {
	s := NewStack(op, e)
	s.Push(op, ev)
	return s
}

func (e VarEvaler) CanEvalFloat(lu Lookuper, funcs *Funcs) bool {
	return lu.CanEvalFloat(string(e))
}

//

type FuncEvaler struct {
	Name string
	Args []Evaler
}

func (e FuncEvaler) Eval(lu Lookuper, funcs *Funcs) (interface{}, error) {
	return funcs.Eval(e.Name, lu, e.Args)
}

func (e FuncEvaler) Push(op Op, ev Evaler) Evaler {
	s := NewStack(op, e)
	s.Push(op, ev)
	return s
}
func (e FuncEvaler) CanEvalFloat(lu Lookuper, funcs *Funcs) bool {
	return funcs.CanEvalFloat(e.Name)
}
