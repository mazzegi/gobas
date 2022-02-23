package gobas

import "strings"

func IsString(varName string) bool {
	return strings.HasSuffix(varName, "$")
}

const StmtSep = ':'

type Stmt interface{}

type ArrayDef struct {
	Var        string
	Dimensions []Expr
}

type DATA struct {
	Consts []string
}

type DEF struct {
	Name string
	Expr Expr
}

type DIM struct {
	Arrays []ArrayDef
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
	Msg       string
	Semicolon bool
	Vars      []string
}

type LET struct {
	Var  string
	Expr Expr
}

type ASSIGN struct {
	Var  string
	Expr Expr
}

type ASSIGN_ARRAY struct {
	Array ArrayDef
	Expr  Expr
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
	Raw   string
	Items []printItem
}

type READ struct {
	Vars []string
}

type REM struct {
	What string
}

type RESTORE struct {
}

type RETURN struct {
}

type STOP struct {
}
