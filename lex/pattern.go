package lex

import (
	"strings"

	"github.com/pkg/errors"
)

//type Whitespace struct{}

type Params map[string]interface{}

type Pattern struct {
	Targets []interface{}
}

func (p *Pattern) AppendTarget(t interface{}) error {
	if _, ok := t.(MatchTarget); ok {
		if len(p.Targets) > 0 {
			if _, ok := p.Targets[len(p.Targets)-1].(MatchTarget); ok {
				return errors.Errorf("a match-target cannot immediately follwo a match-target")
			}
		}
	}
	p.Targets = append(p.Targets, t)
	return nil
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
			err = p.AppendTarget(mt)
			if err != nil {
				return err
			}
		} else {
			str := strings.TrimSpace(curr)
			if str != "" {
				err := p.AppendTarget(str)
				if err != nil {
					return err
				}
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
				//p.AppendTarget(Whitespace{})
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

func (p *Pattern) Eval(s string) (Params, error) {
	s = strings.TrimSpace(s)
	var pos int

	eatWhite := func() int {
		var eaten int
		for {
			if s[pos] != ' ' {
				return eaten
			}
			pos++
			if pos >= len(s) {
				return eaten
			}
			eaten++
		}
	}

	ps := Params{}
	for i, t := range p.Targets {
		if pos >= len(s) {
			return Params{}, errors.Errorf("EOF")
		}
		switch t := t.(type) {
		// case Whitespace:
		// 	wsEaten := eatWhite()
		// 	if wsEaten == 0 {
		// 		return Params{}, errors.Errorf("no whitespace where expected")
		// 	}
		case string:
			eatWhite()
			if !strings.HasPrefix(s[pos:], t) {
				return Params{}, errors.Errorf("no match for string %q", t)
			}
			pos += len(t)
		case MatchTarget:
			eatWhite()
			var smt string
			if i == len(p.Targets)-1 {
				smt = s[pos:]
			} else {
				//peek next string
				next, ok := p.Targets[i+1].(string)
				if !ok {
					return Params{}, errors.Errorf("next is not a string")
				}
				nextIdx := strings.Index(s[pos:], next)
				if nextIdx < 0 {
					return Params{}, errors.Errorf("no match for next %q", next)
				}
				smt = s[pos : pos+nextIdx]
			}
			smt = strings.TrimSpace(smt)

			v, err := t.Eval(smt)
			if err != nil {
				return Params{}, err
			}
			ps[t.Name] = v
			pos += len(smt)
		}
	}

	return ps, nil
}
