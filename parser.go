package gobas

import (
	"fmt"
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
	p.lexer.MustAdd(KeyFOR_STEP, "FOR {ivar:string}={iexpr:string} TO {toexpr:string} STEP {stepexpr:string}")
	p.lexer.MustAdd(KeyFOR, "FOR {ivar:string}={iexpr:string} TO {toexpr:string}")
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
	p.lexer.MustAdd(KeyPRINT, "PRINT{expr:string}")
	p.lexer.MustAdd(KeyPRINT_EMPTY, "PRINT")
	p.lexer.MustAdd(KeyREAD, "READ {expr:string}")
	p.lexer.MustAdd(KeyREM, "REM{expr:string}")
	p.lexer.MustAdd(KeyREM_EMPTY, "REM")
	p.lexer.MustAdd(KeyRESTORE, "RESTORE")
	p.lexer.MustAdd(KeyRETURN, "RETURN")
	p.lexer.MustAdd(KeySTOP, "STOP")
	p.lexer.MustAdd(KeyASSIGN, "{var:string}={expr:string}")
}

func (p *Parser) parseLine(rl rawLine) ([]Stmt, error) {
	if strings.HasPrefix(strings.TrimSpace(rl.text), "REM") {
		return []Stmt{
			REM{What: rl.text},
		}, nil
	}

	stmts, err := p.parseStmts(rl.text)
	if err != nil {
		return nil, errors.Wrapf(err, "in line %d (src = %d)", rl.num, rl.sourceLine)
	}
	return stmts, nil
}

func (p *Parser) parseStmts(s string) ([]Stmt, error) {
	var stmts []Stmt
	stmtsRaw := splitOutsideQuotes(s, StmtSep)
	for _, stmtRaw := range stmtsRaw {
		stmtRaw = trimWhite(stmtRaw)
		if stmtRaw == "" {
			continue
		}
		stmt, err := p.parseStmt(stmtRaw)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) parseStmt(stmtRaw string) (Stmt, error) {
	ps, key, err := p.lexer.Eval(stmtRaw)
	if err != nil {
		return nil, errors.Wrapf(err, "eval stmt %q", stmtRaw)
	}
	fmt.Printf("%q: %s\n", key, ps.Format())
	return nil, nil
}
