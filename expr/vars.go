package expr

import "github.com/pkg/errors"

func NewVars() *Vars {
	return &Vars{
		vars: map[string]interface{}{},
	}
}

type Vars struct {
	vars map[string]interface{}
}

func (vs *Vars) Add(name string, value interface{}) {
	vs.vars[name] = value
}

func (vs *Vars) LookupVar(name string) (interface{}, error) {
	v, ok := vs.vars[name]
	if !ok {
		return nil, errors.Errorf("no such var %q", name)
	}
	return v, nil
}

func (vs *Vars) CanEvalFloat(name string) bool {
	v, ok := vs.vars[name]
	if !ok {
		return false
	}
	return canConvertToFloat(v)
}
