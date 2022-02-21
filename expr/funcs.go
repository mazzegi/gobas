package expr

import (
	"reflect"

	"github.com/pkg/errors"
)

type Func func([]interface{}) (interface{}, error)
type FloatFunc func([]interface{}) (float64, error)

func scanArg(v interface{}, arg interface{}) error {
	rvarg := reflect.ValueOf(arg)
	if rvarg.Kind() != reflect.Pointer {
		return errors.Errorf("cannot scan into a non-pointer")
	}
	argelem := rvarg.Elem()
	if !argelem.CanSet() {
		return errors.Errorf("cannot set %s", argelem.Type().String())
	}

	rv := reflect.ValueOf(v)
	if !rv.CanConvert(argelem.Type()) {
		return errors.Errorf("cannot convert %T to %s", v, argelem.Type().String())
	}

	crv := rv.Convert(argelem.Type())
	argelem.Set(crv)

	//reflect.ValueOf(v).CanConvert()

	return nil
}

func ScanArgs(vs []interface{}, args ...interface{}) error {
	if len(vs) != len(args) {
		return errors.Errorf("expect %d args, got %d", len(args), len(vs))
	}
	for i, v := range vs {
		arg := args[i]
		err := scanArg(v, arg)
		if err != nil {
			return errors.Wrapf(err, "scan arg %d", i)
		}
	}

	return nil
}

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
