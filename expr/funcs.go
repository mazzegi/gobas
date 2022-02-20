package expr

import "github.com/pkg/errors"

type Func func([]interface{}) (interface{}, error)

func NewFuncs() *Funcs {
	return &Funcs{
		items: map[string]Func{},
	}
}

type Funcs struct {
	items map[string]Func
}

func (fs *Funcs) Add(name string, fnc Func) {
	fs.items[name] = fnc
}

func (fs *Funcs) Eval(name string, lu Lookuper, evs []Evaler) (interface{}, error) {
	fnc, ok := fs.items[name]
	if !ok {
		return 0, errors.Errorf("no such func %q", name)
	}
	var vs []interface{}
	for _, ev := range evs {
		v, err := ev.Eval(lu, fs)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return fnc(vs)
}

func (fs *Funcs) CanEvalFloat(name string) bool {
	//TODO: check
	return true
}

func (fs *Funcs) Contains(name string) bool {
	_, ok := fs.items[name]
	return ok
}
