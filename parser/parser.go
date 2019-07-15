package parser

import (
	"github.com/pkg/errors"
	"lox/scanner"
)

/*
program     -> declaration* EOF ;
declaration -> funDec | classDec | varDedc | statement ;
statement   -> exprStatement | printStatemt | block | ifStatement | whileStatement | jumpStatement | returnStmt;
varDec      -> "var" IDENTIFIER ("=" expresssion) ;
funDec      -> "fun" IDENTIFIER "(" (IDENTIFIER ",")* ")" funExpr ;
expression  -> assign
assign      -> IDENTIFIER "=" assign * | IDENTIFIER assignOp expression | logicOr ;
assignOp    -> "+=" | "-=" | "*=" | "/=" ;
logicOr     -> logicAnd ("or" logicAnd)* ;
logicAnd    -> equality ("and" equality)* ;
equality    -> tern (("!=" | "==") tern)* ;
tern        -> comparisson "?" comparisson ":" tern ;
comparisson -> addition ((">" | "<" | "<=" | ">=") addition)* ;
addition    -> mult (("+" | "-") mult)* ;
mult        -> unary ("*" | "/") unary)* ;
unary       -> "+" | "-" unary | call;
call        -> primary |
primary     -> IDENTIFIER | funExpr | NUMBER | STRING | "(" expression ")" ;
 */

type Parser struct {
	Program *ASTNode
	Tokens []scanner.Token
	Current int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{nil, tokens, 0}
}

func (p *Parser) parse() ([]Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		dec, _ := p.declaration()
		statements = append(statements, dec)
	}
	return statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(scanner.VAR) {
		name, _ := p.consume(scanner.IDENTIFIER, "expected identifier")
		var init Expr = nil
		if p.match(scanner.EQUAL) {
			init, _ = p.assignment()
		}
		return &VarStmt{*name, init}, nil
	}
	statements, _ := p.statement()
	return statements, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(scanner.IF) {
		stmt, _ := p.ifStmt()
		return &stmt, nil
	}
	if p.match(scanner.WHILE) {
		stmt, _ := p.whileStmt()
		return &stmt, nil
	}
	if p.match(scanner.LB) {
		stmt, _ := p.block()
		return &stmt, nil
	}
	stmt, _ := p.exprStmt()
	return &stmt, nil
}

func (p *Parser) exprStmt() (ExprStmt, error) {
	expr, _ := p.assignment()
	p.consume(scanner.SEMICOLON, "expected ';' after expression")
	return ExprStmt{expr}, nil
}

func (p *Parser) ifStmt() (IfStmt, error) {
	p.consume(scanner.LP, "expect '(' after 'if'")
	cond, _ := p.assignment()
	p.consume(scanner.RP, "expect ')' after 'if' condition")
	thenBranch, _ := p.block()
	var elseBranch Stmt = nil
	if p.match(scanner.ELSE) {
		elseBranch, _ = p.statement()
	}
	return IfStmt{cond, &thenBranch, elseBranch}, nil
}

func (p *Parser) whileStmt() (WhileStmt, error) {
	p.consume(scanner.LP, "expect '(' after 'if'")
	cond, _ := p.assignment()
	p.consume(scanner.RP, "expect ')' after 'if' condition")
	body, _ := p.statement()
	return WhileStmt{cond, body}, nil
}

func (p *Parser) block() (BlockStmt, error) {
	var statements []Stmt
	for !p.isAtEnd() && !p.check(scanner.RB) {
		stmt, _ := p.declaration()
		statements = append(statements, stmt)
	}
	return BlockStmt{statements}, nil
}

func (p *Parser) assignment() (Expr, error) {
	expr, _ := p.logicalOr()
	for p.match(scanner.EQUAL) {
		op := p.Tokens[p.Current-1]
		expr, _ := p.assignment()
		return &AssignExpr{op, expr}, nil
	}
	return expr, nil
}

func (p *Parser) logicalOr() (Expr, error) {
	lv, _ := p.logicalAnd()
	for p.match(scanner.OR) {
		op := p.Tokens[p.Current-1]
		rv, _ := p.logicalAnd()
		lv = &LogicalExpr{lv, op, rv}
	}
	return lv, nil
}

func (p *Parser) logicalAnd() (Expr, error) {
	lv, _ := p.equality()
	for p.match(scanner.AND) {
		op := p.Tokens[p.Current-1]
		rv, _ := p.equality()
		lv = &BinExpr{lv, op, rv}
	}
	return lv, nil
}

func (p *Parser) equality() (Expr, error) {
	lv, _ := p.comparison()
	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		op := p.Tokens[p.Current-1]
		rv, _ := p.comparison()
		lv = &BinExpr{lv, op, rv}
	}
	return lv, nil
}

func (p *Parser) comparison() (Expr, error) {
	lv, _ := p.addition()
	for p.match(scanner.GRATER, scanner.LESS,
		scanner.LESS_EQUAL, scanner.GREATER_EQUAL) {
		op := p.Tokens[p.Current-1]
		rv, _ := p.addition()
		lv = &BinExpr{lv, op, rv}
	}
	return lv, nil
}

func (p *Parser) addition() (Expr, error) {
	lv, _ := p.multiplication()
	for p.match(scanner.PLUS, scanner.MINUS) {
		op := p.Tokens[p.Current-1]
		rv, _ := p.multiplication()
		lv = &BinExpr{lv, op, rv}
	}
	return lv, nil
}

func (p *Parser) multiplication() (Expr, error) {
	lv, _ := p.unary()
	for p.match(scanner.STAR, scanner.SLASH) {
		op := p.Tokens[p.Current-1]
		rv, _ := p.unary()
		lv = &BinExpr{lv, op, rv}
	}
	return lv, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(scanner.BANG, scanner.MINUS, scanner.PLUS) {
		op := p.Tokens[p.Current-1]
		rv, _ := p.assignment()
		return &UnExpr{op, rv}, nil
	}
	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, _ := p.primary()
	var args []Expr
	for p.match(scanner.LP) {
		args, _ = p.finishCall()
	}
	return &CallExpr{expr, args}, nil
}

func (p *Parser) finishCall() ([]Expr, error) {
	var args []Expr
	for !p.match(scanner.RP) {
		cond := true
		for cond {
			tmp, _ := p.assignment()
			args = append(args, tmp)
			cond = p.match(scanner.COMMA)
		}
	}
	return args, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(scanner.FALSE) {
		return &LiteralExpr{false}, nil
	}
	if p.match(scanner.TRUE) {
		return &LiteralExpr{true}, nil
	}
	if p.match(scanner.NIL) {
		return &LiteralExpr{nil}, nil
	}
	if p.match(scanner.NIL) {
		return &LiteralExpr{nil}, nil
	}
	if p.match(scanner.NUMBER, scanner.STRING) {
		return &LiteralExpr{p.Tokens[p.Current-1].Literal}, nil
	}
	if p.match(scanner.IDENTIFIER) {
		return &VarExpr{p.Tokens[p.Current-1]}, nil
	}
	return nil, nil
}

func (p *Parser) consume(tokenType scanner.TokenType, message string) (*scanner.Token, error) {
	token := p.advance()
	if token.Type != tokenType {
		return nil, errors.New(message)
	}

}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.Tokens[p.Current-1]
}

func (p *Parser) match(types... scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	for _, tp := range types {
		if p.Tokens[p.Current].Type == tp {
			p.Current++
			return true
		}
	}
	return false
}

func (p *Parser) check(types... scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	for _, tp := range types {
		if p.Tokens[p.Current].Type == tp {
			return true
		}
	}
	return false
}

func (p *Parser) isAtEnd() bool {
	return p.Current >= len(p.Tokens)
}
