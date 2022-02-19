package gobas

const StmtSep = ";"

type Stmt interface{}

type DATA struct {
	//TODO
}

type DEF struct {
	Name string
	Expr Expr
}

type DIM struct {
	Size     uint32
	IsString bool
}

type END struct{}

type FOR struct {
	Var     string
	Initial Expr
	To      Expr
	Step    Expr
}

type GOSUB struct {
	Line uint32
}

type GOTO struct {
	Line uint32
}

type IFLN struct {
	Expr Expr
	Line uint32
}

type IFELSELN struct {
	Expr     Expr
	Line     uint32
	ElseLine uint32
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
	Lines []uint32
}

type ONGOTO struct {
	Expr  Expr
	Lines []uint32
}

type PRINT struct {
	Exprs []Expr
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
