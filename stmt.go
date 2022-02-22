package gobas

import "strings"

func IsString(varName string) bool {
	return strings.HasSuffix(varName, "$")
}

const StmtSep = ':'

type Stmt interface{}

type Array struct {
	Var        string
	Dimensions []int
}

type DATA struct {
	//TODO
	Expr string
}

type DEF struct {
	Name string
	Expr Expr
}

type DIM struct {
	Arrays []Array
}

type END struct{}

type FOR struct {
	Var     string
	Initial Expr
	To      Expr
	Step    Expr
}

type GOSUB struct {
	Line int
}

type GOTO struct {
	Line int
}

type IFLN struct {
	Expr Expr
	Line int
}

type IFELSELN struct {
	Expr     Expr
	Line     int
	ElseLine int
}

type IFSTMT struct {
	Expr  Expr
	Stmts []Stmt
}

type IFELSESTMT struct {
	Expr      Expr
	Stmts     []Stmt
	ElseStmts []Stmt
}

type INPUT struct {
	//TODO
}

type LET struct {
	Var  string
	Expr Expr
}

type ASSIGN struct {
	Var  string
	Expr Expr
}

type NEXT struct {
	Var string
}

type ONGOSUB struct {
	Expr  Expr
	Lines []int
}

type ONGOTO struct {
	Expr  Expr
	Lines []int
}

type PRINT struct {
	//Exprs []Expr
	Raw   string
	Items []printItem
}

type READ struct {
	//TODO
}

type REM struct {
	What string
}

type RESTORE struct {
	//TODO
}

type RETURN struct {
}

type STOP struct {
}
