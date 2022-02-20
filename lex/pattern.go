package lex

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Params map[string]interface{}

func (ps Params) Format() string {
	var sl []string
	for k, v := range ps {
		sl = append(sl, fmt.Sprintf("%q =[%v]", k, v))
	}
	return strings.Join(sl, ", ")
}

type Pattern struct {
	Prefix  string
	Name    string
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

	// if first target is a string, take it as prefix
	if st, ok := t.(string); ok && len(p.Targets) == 0 {
		p.Prefix = st
	}

	p.Targets = append(p.Targets, t)
	return nil
}

func ParsePatternWithName(name string, s string) (*Pattern, error) {
	p, err := ParsePattern(s)
	if err != nil {
		return nil, err
	}
	p.Name = name
	return p, nil
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
		case string:
			eatWhite()
			if !strings.HasPrefix(s[pos:], t) {
				return Params{}, errors.Errorf("no match for string %q", t)
			}
			pos += len(t)
		case MatchTarget:
			eatWhite()
			var mts string
			if i == len(p.Targets)-1 {
				mts = s[pos:]
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
				mts = s[pos : pos+nextIdx]
			}
			mts = strings.TrimSpace(mts)

			v, err := t.Eval(mts)
			if err != nil {
				return Params{}, err
			}
			ps[t.Name] = v
			pos += len(mts)
		}
	}

	return ps, nil
}
