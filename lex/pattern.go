package lex

import (
	"strings"

	"github.com/pkg/errors"
)

type Whitespace struct{}

type Pattern struct {
	Targets []interface{}
}

func (p *Pattern) AppendTarget(t interface{}) {
	if _, ok := t.(Whitespace); ok {
		if len(p.Targets) > 0 {
			if _, ok := p.Targets[len(p.Targets)-1].(Whitespace); ok {
				return
			}
		}
	}
	p.Targets = append(p.Targets, t)
}

// ParsePattern parses patterns like "ON {expr:string} GOSUB {lines:[]int}"
func ParsePattern(s string) (*Pattern, error) {
	p := &Pattern{}
	var curr string
	var inCurly bool

	flushCurrent := func() error {
		if inCurly {
			mt, err := ParseMatchTarget(curr)
			if err != nil {
				return errors.Wrapf(err, "parse match-target %q", curr)
			}
			p.AppendTarget(mt)
		} else {
			str := strings.TrimSpace(curr)
			if str != "" {
				p.AppendTarget(str)
			}
		}
		curr = ""
		return nil
	}

	for _, r := range s {
		if r == '{' {
			err := flushCurrent()
			if err != nil {
				return nil, err
			}
			if inCurly {
				return nil, errors.Errorf("found { in match-target-expression")
			}
			inCurly = true
			continue
		}
		if r == '}' {
			if !inCurly {
				return nil, errors.Errorf("found } outside match-target-expression")
			}
			err := flushCurrent()
			if err != nil {
				return nil, err
			}
			inCurly = false
			continue
		}
		if inCurly {
			curr += string(r)
		} else {
			if r == ' ' {
				err := flushCurrent()
				if err != nil {
					return nil, err
				}
				p.AppendTarget(Whitespace{})
			} else {
				curr += string(r)
			}
		}
	}
	err := flushCurrent()
	if err != nil {
		return nil, err
	}
	return p, nil
}
