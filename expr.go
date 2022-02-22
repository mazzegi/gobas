package gobas

import (
	"github.com/mazzegi/gobas/expr"
	"github.com/pkg/errors"
)

func parseExpression(s string) (Expr, error) {
	if s == "" {
		return Expr{}, nil
	}

	stack, err := expr.NewParser(s).Parse()
	if err != nil {
		return Expr{}, err
	}

	ex := Expr{
		Raw:   s,
		Stack: stack,
	}
	return ex, nil
}

func mustParseExpression(s string) Expr {
	ex, err := parseExpression(s)
	if err != nil {
		panic(err)
	}
	return ex
}

func parseExpressions(sl []string) ([]Expr, error) {
	var es []Expr
	for _, s := range sl {
		e, err := parseExpression(s)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-expression %q", s)
		}
		es = append(es, e)
	}
	return es, nil
}

func mustParseExpressions(sl []string) []Expr {
	es, err := parseExpressions(sl)
	if err != nil {
		panic(err)
	}
	return es
}

type Expr struct {
	Raw   string
	Stack *expr.Stack
}
