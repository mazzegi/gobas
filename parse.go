package gobas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func ParseFile(fileName string) ([]Stmt, error) {
	rls, err := rawReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var stmts []Stmt
	for _, rl := range rls {
		lineStmts, err := parseLine(rl)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, lineStmts...)
	}
	return stmts, nil
}

type parseFunc func(string) (Stmt, error)

var parseFuncs map[string]parseFunc

func init() {
	parseFuncs = map[string]parseFunc{
		"DATA":    parseDATA,
		"DEF":     parseDEF,
		"DIM":     parseDIM,
		"END":     parseEND,
		"FOR":     parseFOR,
		"GOSUB":   parseGOSUB,
		"GOTO":    parseGOTO,
		"IF":      parseIF,
		"LET":     parseLET,
		"NEXT":    parseNEXT,
		"ON":      parseON,
		"READ":    parseREAD,
		"REM":     parseREM,
		"RESTORE": parseRESTORE,
		"RETURN":  parseRETURN,
		"STOP":    parseSTOP,
	}
}

func parseLine(rl rawLine) ([]Stmt, error) {
	return parseStmts(rl.text)
}

func parseStmts(s string) ([]Stmt, error) {
	var stmts []Stmt
	stmtsRaw := splitOutsideQuotes(s, StmtSep)
	for _, stmtRaw := range stmtsRaw {
		stmt, err := parseStmt(stmtRaw)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func parseStmt(stmtRaw string) (Stmt, error) {
	stmtRaw = trimWhite(stmtRaw)
	if stmtRaw == "" {
		return nil, errors.Errorf("empty statement")
	}

	// PRINT and INPUT may have no space after them
	switch {
	case strings.HasPrefix(stmtRaw, "PRINT"):
		return parsePRINT(stmtRaw)
	case strings.HasPrefix(stmtRaw, "INPUT"):
		return parseINPUT(stmtRaw)
	}

	// check for keywords
	// as stmtRaw is not empty, strings.Fields give always at least on elt
	key := strings.Fields(stmtRaw)[0]
	if pf, ok := parseFuncs[key]; ok {
		data := trimWhite(strings.TrimPrefix(stmtRaw, key))
		return pf(data)
	}

	// no keyword - then maybe an assignment
	return parseASSIGN(stmtRaw)
}

// parse funcs

func parseDATA(s string) (Stmt, error) {
	//TODO: impl
	stmt := DATA{}
	return stmt, nil
}

func parseDEF(s string) (Stmt, error) {
	var name, exprStr string
	_, err := fmt.Sscanf(s, "%s=%s", &name, &exprStr)
	if err != nil {
		return nil, err
	}
	expr, err := parseExpr(exprStr)
	if err != nil {
		return nil, err
	}
	stmt := DEF{
		Name: name,
		Expr: expr,
	}
	return stmt, nil
}

func parsePRINT(s string) (Stmt, error) {
	sl := splitOutsideQuotes(s, ';')
	stmt := PRINT{}
	for _, se := range sl {
		ex, err := parseExpr(se)
		if err != nil {
			return nil, err
		}
		stmt.Exprs = append(stmt.Exprs, ex)
	}
	return stmt, nil
}

func parseINPUT(s string) (Stmt, error) {
	//TODO
	stmt := INPUT{}
	return stmt, nil
}

func parseASSIGN(s string) (Stmt, error) {
	var varName string
	var exprStr string
	_, err := fmt.Sscanf(s, "%s=%s", &varName, &exprStr)
	if err != nil {
		return nil, err
	}
	expr, err := parseExpr(exprStr)
	if err != nil {
		return nil, err
	}

	return ASSIGN{
		Var:  varName,
		Expr: expr,
	}, nil
}

func parseDIM(s string) (Stmt, error) {
	stmt := DIM{}
	sl := splitOutsideQuotes(s, ',')
	for _, sa := range sl {
		var name string
		var dims string
		_, err := fmt.Sscanf(sa, "%s(%s)", &name, &dims)
		if err != nil {
			return nil, err
		}
		ns, err := parseInts[int](dims, ',')
		if err != nil {
			return nil, err
		}
		stmt.Arrays = append(stmt.Arrays, Array{
			Var:        name,
			Dimensions: ns,
		})
	}
	return stmt, nil
}

func parseEND(s string) (Stmt, error) {
	stmt := END{}
	return stmt, nil
}

func parseFOR(s string) (Stmt, error) {
	var varName string
	var initExprStr string
	var toExprStr string
	var stepStr string
	var err error
	if strings.Contains(s, " STEP ") {
		_, err = fmt.Sscanf(s, "%s=%s TO %s STEP %s", &varName, &initExprStr, &toExprStr, &stepStr)
	} else {
		_, err = fmt.Sscanf(s, "%s=%s TO %s", &varName, &initExprStr, &toExprStr)
	}
	if err != nil {
		return nil, err
	}
	if stepStr == "" {
		stepStr = "1"
	}
	initExpr, err := parseExpr(initExprStr)
	if err != nil {
		return nil, err
	}
	toExpr, err := parseExpr(toExprStr)
	if err != nil {
		return nil, err
	}
	stepExpr, err := parseExpr(stepStr)
	if err != nil {
		return nil, err
	}

	stmt := FOR{
		Var:     varName,
		Initial: initExpr,
		To:      toExpr,
		Step:    stepExpr,
	}
	return stmt, nil
}

func parseGOSUB(s string) (Stmt, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil, err
	}
	stmt := GOSUB{
		Line: uint32(n),
	}
	return stmt, nil
}

func parseGOTO(s string) (Stmt, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil, err
	}
	stmt := GOTO{
		Line: uint32(n),
	}
	return stmt, nil
}

func parseIF(s string) (Stmt, error) {
	var condExprStr string
	var rest string
	_, err := fmt.Sscanf(s, "%s THEN %s", &condExprStr, &rest)
	if err != nil {
		return nil, err
	}
	condExpr, err := parseExpr(condExprStr)
	if err != nil {
		return nil, err
	}

	if strings.Contains(rest, " ELSE ") {
		var ifDo string
		var elseDo string
		_, err := fmt.Sscanf(rest, "%s ELSE %s", &ifDo, &elseDo)
		if err != nil {
			return nil, err
		}
		// is ifDo a (line)number
		if ifDoLine, err := strconv.ParseUint(ifDo, 10, 32); err == nil {
			// is elseDo also ?
			if elseDoLine, err := strconv.ParseUint(elseDo, 10, 32); err == nil {
				return IFELSELN{
					Expr:     condExpr,
					Line:     uint32(ifDoLine),
					ElseLine: uint32(elseDoLine),
				}, nil
			}
			return nil, errors.Errorf("in if else statements all branches must be either lines or statements")
		} else {
			ifDoStmts, err := parseStmts(ifDo)
			if err != nil {
				return nil, err
			}
			elseDoStmts, err := parseStmts(elseDo)
			if err != nil {
				return nil, err
			}
			return IFELSESTMT{
				Expr:      condExpr,
				Stmts:     ifDoStmts,
				ElseStmts: elseDoStmts,
			}, nil
		}
	} else {
		// is it a (line)number
		if n, err := strconv.ParseUint(rest, 10, 32); err == nil {
			return IFLN{
				Expr: condExpr,
				Line: uint32(n),
			}, nil
		} else {
			stmts, err := parseStmts(rest)
			if err != nil {
				return nil, err
			}
			return IFSTMT{
				Expr:  condExpr,
				Stmts: stmts,
			}, nil
		}
	}
}

func parseLET(s string) (Stmt, error) {
	var varName string
	var exprStr string
	_, err := fmt.Sscanf(s, "%s=%s", &varName, &exprStr)
	if err != nil {
		return nil, err
	}
	expr, err := parseExpr(exprStr)
	if err != nil {
		return nil, err
	}

	return LET{
		Var:  varName,
		Expr: expr,
	}, nil
}

func parseNEXT(s string) (Stmt, error) {
	varName := trimWhite(s)
	if varName == "" {
		return nil, errors.Errorf("var name is empty")
	}
	return NEXT{
		Var: varName,
	}, nil
}

func parseON(s string) (Stmt, error) {
	var exprStr string
	var linesStr string
	if strings.Contains(s, " GOSUB ") {
		_, err := fmt.Sscanf(s, "%s GOSUB %s", &exprStr, &linesStr)
		if err != nil {
			return nil, err
		}
		expr, err := parseExpr(exprStr)
		if err != nil {
			return nil, err
		}
		lines, err := parseInts[uint32](linesStr, ',')
		if err != nil {
			return nil, err
		}
		return ONGOSUB{
			Expr:  expr,
			Lines: lines,
		}, nil
	} else if strings.Contains(s, " GOTO ") {
		_, err := fmt.Sscanf(s, "%s GOTO %s", &exprStr, &linesStr)
		if err != nil {
			return nil, err
		}
		expr, err := parseExpr(exprStr)
		if err != nil {
			return nil, err
		}
		lines, err := parseInts[uint32](linesStr, ',')
		if err != nil {
			return nil, err
		}
		return ONGOSUB{
			Expr:  expr,
			Lines: lines,
		}, nil
	} else {
		return nil, errors.Errorf("invalid ON statement %q", s)
	}
}

func parseREAD(s string) (Stmt, error) {
	//TODO
	stmt := READ{}
	return stmt, nil
}
func parseREM(s string) (Stmt, error) {
	stmt := REM{
		What: s,
	}
	return stmt, nil
}
func parseRESTORE(s string) (Stmt, error) {
	//TODO
	stmt := RESTORE{}
	return stmt, nil
}
func parseRETURN(s string) (Stmt, error) {
	//TODO
	stmt := RETURN{}
	return stmt, nil
}
func parseSTOP(s string) (Stmt, error) {
	//TODO
	stmt := STOP{}
	return stmt, nil
}
