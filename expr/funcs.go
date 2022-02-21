package expr

import "github.com/pkg/errors"

type Func func([]interface{}) (interface{}, error)
type FloatFunc func([]interface{}) (float64, error)

func NewFuncs() *Funcs {
	return &Funcs{
		funcs:      map[string]Func{},
		floatFuncs: map[string]FloatFunc{},
	}
}

type Funcs struct {
	funcs      map[string]Func
	floatFuncs map[string]FloatFunc
}

func (fs *Funcs) AddFunc(name string, fnc Func) {
	fs.funcs[name] = fnc
}

func (fs *Funcs) AddFloatFunc(name string, fnc FloatFunc) {
	fs.floatFuncs[name] = fnc
}

func (fs *Funcs) Eval(name string, lu Lookuper, evs []Evaler) (interface{}, error) {
	var vs []interface{}
	for _, ev := range evs {
		v, err := ev.Eval(lu, fs)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}

	if ffnc, ok := fs.floatFuncs[name]; ok {
		return ffnc(vs)
	}
	if fnc, ok := fs.funcs[name]; ok {
		return fnc(vs)
	}
	return nil, errors.Errorf("no such func %q", name)
}

func (fs *Funcs) CanEvalFloat(name string) bool {
	if _, ok := fs.floatFuncs[name]; ok {
		return true
	}
	return false
}

func (fs *Funcs) Contains(name string) bool {
	if _, ok := fs.floatFuncs[name]; ok {
		return true
	}
	if _, ok := fs.funcs[name]; ok {
		return true
	}
	return false
}
