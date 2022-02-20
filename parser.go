package gobas

import (
	"strings"

	"github.com/mazzegi/gobas/lex"
	"github.com/pkg/errors"
)

func NewParser() *Parser {
	p := &Parser{}
	p.init()
	return p
}

type Parser struct {
	lexer *lex.Set
}

func (p *Parser) ParseFile(fileName string) ([]Stmt, error) {
	rls, err := rawReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var stmts []Stmt
	for _, rl := range rls {
		lineStmts, err := p.parseLine(rl)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, lineStmts...)
	}
	return stmts, nil
}

// private stuff

const (
	KeyDATA        = "DATA"
	KeyDEF         = "DEF"
	KeyDIM         = "DIM"
	KeyEND         = "END"
	KeyFOR         = "FOR"
	KeyFOR_STEP    = "FOR_STEP"
	KeyGOSUB       = "GOSUB"
	KeyGOTO        = "GOTO"
	KeyIFLN        = "IFLN"
	KeyIFELSELN    = "IFELSELN"
	KeyIFSTMT      = "IFSTMT"
	KeyIFELSESTMT  = "IFELSESTMT"
	KeyINPUT       = "INPUT"
	KeyLET         = "LET"
	KeyNEXT        = "NEXT"
	KeyNEXT_EMPTY  = "NEXT_EMPTY"
	KeyON_GOSUB    = "ON_GOSUB"
	KeyON_GOTO     = "ON_GOTO"
	KeyPRINT       = "PRINT"
	KeyPRINT_EMPTY = "PRINT_EMPTY"
	KeyREAD        = "READ"
	KeyREM         = "REM"
	KeyREM_EMPTY   = "REM_EMPTY"
	KeyRESTORE     = "RESTORE"
	KeyRETURN      = "RETURN"
	KeySTOP        = "STOP"
	KeyASSIGN      = "ASSIGN"
)

func (p *Parser) init() {
	p.lexer = lex.NewSet()
	p.lexer.MustAdd(KeyDATA, "DATA {expr:string}")
	p.lexer.MustAdd(KeyDEF, "DEF {fnc:string}={expr:string}")
	p.lexer.MustAdd(KeyDIM, "DIM {arrayexprs:[]string?sep=,}")
	p.lexer.MustAdd(KeyEND, "END")
	p.lexer.MustAdd(KeyFOR_STEP, "FOR {var:string}={iexpr:string} TO {toexpr:string} STEP {stepexpr:string}")
	p.lexer.MustAdd(KeyFOR, "FOR {var:string}={iexpr:string} TO {toexpr:string}")
	p.lexer.MustAdd(KeyGOSUB, "GOSUB {line:int}")
	p.lexer.MustAdd(KeyGOTO, "GOTO {line:int}")
	p.lexer.MustAdd(KeyIFELSELN, "IF {condexpr:string} THEN {line:int} ELSE {elseline:int}")
	p.lexer.MustAdd(KeyIFELSESTMT, "IF {condexpr:string} THEN {stmts:string} ELSE {elsestmts:string}")
	p.lexer.MustAdd(KeyIFLN, "IF {condexpr:string} THEN {line:int}")
	p.lexer.MustAdd(KeyIFSTMT, "IF {condexpr:string} THEN {stmts:string}")
	p.lexer.MustAdd(KeyINPUT, "INPUT{expr:string}")
	p.lexer.MustAdd(KeyLET, "LET {var:string}={expr:string}")
	p.lexer.MustAdd(KeyNEXT, "NEXT {var:string}")
	p.lexer.MustAdd(KeyNEXT_EMPTY, "NEXT")
	p.lexer.MustAdd(KeyON_GOSUB, "ON {expr:string} GOSUB {lines:[]int}")
	p.lexer.MustAdd(KeyON_GOTO, "ON {expr:string} GOTO {lines:[]int}")
	p.lexer.MustAdd(KeyPRINT, "PRINT{exprs:[]string?sep=;}")
	p.lexer.MustAdd(KeyPRINT_EMPTY, "PRINT")
	p.lexer.MustAdd(KeyREAD, "READ {expr:string}")
	p.lexer.MustAdd(KeyREM, "REM{expr:string}")
	p.lexer.MustAdd(KeyREM_EMPTY, "REM")
	p.lexer.MustAdd(KeyRESTORE, "RESTORE")
	p.lexer.MustAdd(KeyRETURN, "RETURN")
	p.lexer.MustAdd(KeySTOP, "STOP")
	p.lexer.MustAdd(KeyASSIGN, "{var:string}={expr:string}")
}

func (p *Parser) parseLine(rl rawLine) (stmts []Stmt, err error) {
	if strings.HasPrefix(strings.TrimSpace(rl.text), "REM") {
		return []Stmt{
			REM{What: rl.text},
		}, nil
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("in line %d (src = %d): %v", rl.num, rl.sourceLine+1, r)
		}
	}()
	stmts = p.mustParseStmts(rl.text)
	return
}

func (p *Parser) mustParseStmts(s string) []Stmt {
	var stmts []Stmt
	stmtsRaw := splitOutsideQuotes(s, StmtSep)
	for _, stmtRaw := range stmtsRaw {
		stmtRaw = trimWhite(stmtRaw)
		if stmtRaw == "" {
			continue
		}
		stmt := p.mustParseStmt(stmtRaw)
		stmts = append(stmts, stmt)
	}
	return stmts
}

func (p *Parser) mustParseStmt(stmtRaw string) Stmt {
	ps, key, err := p.lexer.Eval(stmtRaw)
	if err != nil {
		panic(errors.Wrapf(err, "eval stmt %q", stmtRaw))
	}
	//fmt.Printf("%q: %s\n", key, ps.Format())

	switch key {
	case KeyDATA:
		return DATA{
			Expr: lex.MustParam[string](ps, "expr"),
		}
	case KeyDEF:
		return DEF{
			Name: lex.MustParam[string](ps, "fnc"),
			Expr: mustParseExpression(lex.MustParam[string](ps, "expr")),
		}
	case KeyDIM:
		return DIM{
			Arrays: mustParseArrays(lex.MustParam[[]string](ps, "arrayexprs")),
		}
	case KeyEND:
		return END{}
	case KeyFOR:
		return FOR{
			Var:     lex.MustParam[string](ps, "var"),
			Initial: mustParseExpression(lex.MustParam[string](ps, "iexpr")),
			To:      mustParseExpression(lex.MustParam[string](ps, "toexpr")),
			Step:    mustParseExpression("1"),
		}
	case KeyFOR_STEP:
		return FOR{
			Var:     lex.MustParam[string](ps, "var"),
			Initial: mustParseExpression(lex.MustParam[string](ps, "iexpr")),
			To:      mustParseExpression(lex.MustParam[string](ps, "toexpr")),
			Step:    mustParseExpression(lex.MustParam[string](ps, "stepexpr")),
		}
	case KeyGOSUB:
		return GOSUB{
			Line: lex.MustParam[int](ps, "line"),
		}
	case KeyGOTO:
		return GOTO{
			Line: lex.MustParam[int](ps, "line"),
		}
	case KeyIFLN:
		return IFLN{
			Expr: mustParseExpression(lex.MustParam[string](ps, "condexpr")),
			Line: lex.MustParam[int](ps, "line"),
		}
	case KeyIFELSELN:
		return IFELSELN{
			Expr:     mustParseExpression(lex.MustParam[string](ps, "condexpr")),
			Line:     lex.MustParam[int](ps, "line"),
			ElseLine: lex.MustParam[int](ps, "elseline"),
		}
	case KeyIFSTMT:
		return IFSTMT{
			Expr:  mustParseExpression(lex.MustParam[string](ps, "condexpr")),
			Stmts: p.mustParseStmts(lex.MustParam[string](ps, "stmts")),
		}
	case KeyIFELSESTMT:
		return IFELSESTMT{
			Expr:      mustParseExpression(lex.MustParam[string](ps, "condexpr")),
			Stmts:     p.mustParseStmts(lex.MustParam[string](ps, "stmts")),
			ElseStmts: p.mustParseStmts(lex.MustParam[string](ps, "elsestmts")),
		}
	case KeyINPUT:
		return INPUT{}
	case KeyLET:
		return LET{
			Var:  lex.MustParam[string](ps, "var"),
			Expr: mustParseExpression(lex.MustParam[string](ps, "expr")),
		}
	case KeyNEXT:
		return NEXT{
			Var: lex.MustParam[string](ps, "var"),
		}
	case KeyNEXT_EMPTY:
		return NEXT{}
	case KeyON_GOSUB:
		return ONGOSUB{
			Expr:  mustParseExpression(lex.MustParam[string](ps, "expr")),
			Lines: lex.MustParam[[]int](ps, "lines"),
		}
	case KeyON_GOTO:
		return ONGOTO{
			Expr:  mustParseExpression(lex.MustParam[string](ps, "expr")),
			Lines: lex.MustParam[[]int](ps, "lines"),
		}
	case KeyPRINT:
		return PRINT{
			Exprs: mustParseExpressions(lex.MustParam[[]string](ps, "exprs")),
		}
	case KeyPRINT_EMPTY:
		return PRINT{}
	case KeyREAD:
		return READ{}
	case KeyRESTORE:
		return RESTORE{}
	case KeyRETURN:
		return RETURN{}
	case KeySTOP:
		return STOP{}
	case KeyASSIGN:
		return ASSIGN{
			Var:  lex.MustParam[string](ps, "var"),
			Expr: mustParseExpression(lex.MustParam[string](ps, "expr")),
		}
	case KeyREM:
		return REM{
			What: lex.MustParam[string](ps, "expr"),
		}
	case KeyREM_EMPTY:
		return REM{}
	default:
		panic(errors.Errorf("unknown key %q", key))
	}
}

func mustParseArrays(sl []string) []Array {
	var as []Array
	return as
}
