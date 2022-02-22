package boolexpr

import (
	"github.com/mazzegi/gobas/expr"
	"github.com/pkg/errors"
)

func NewParser(expr string, funcs *expr.Funcs) *Parser {
	return &Parser{
		expr:  expr,
		pos:   0,
		funcs: funcs,
	}
}

type Parser struct {
	expr  string
	pos   int
	funcs *expr.Funcs
}

func (p *Parser) Parse() (*Stack, error) {
	var stack *Stack

	if stack == nil {
		return nil, errors.Errorf("stack is empty")
	}
	return stack, nil
}
