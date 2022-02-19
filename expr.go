package gobas

func parseExpr(s string) (Expr, error) {
	ex := Expr{Raw: s}
	return ex, nil
}

type Expr struct {
	Raw string
}
