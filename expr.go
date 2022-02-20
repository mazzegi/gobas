package gobas

import "github.com/pkg/errors"

func parseExpression(s string) (Expr, error) {
	ex := Expr{Raw: s}
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
	Raw string
}
