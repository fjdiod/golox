package parser

import (
	"lox/scanner"
)

type Stmt interface {
	ASTNode
	isStmt()
}

type FunStmt struct {
	Name scanner.Token
	Fun FunExpr
}

func (e *FunStmt) isAstNode() {}
func (e *FunStmt) isStmt() {}

type IfStmt struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e *IfStmt) isAstNode() {}
func (e *IfStmt) isStmt() {}

type WhileStmt struct {
	Condition Expr
	Body Stmt
}

func (e *WhileStmt) isAstNode() {}
func (e *WhileStmt) isStmt() {}

type BlockStmt struct {
	Statements []Stmt
}

func (e *BlockStmt) isAstNode() {}
func (e *BlockStmt) isStmt() {}

type ReturnStmt struct {
	Keyword scanner.Token
	Value Expr
}

func (e *ReturnStmt) isAstNode() {}
func (e *ReturnStmt) isStmt() {}

type JumpStmt struct {
	Keyword scanner.Token
}

func (e *JumpStmt) isAstNode() {}
func (e *JumpStmt) isStmt() {}

type PrintStmt struct {
	Keyword scanner.Token
	Value Expr
}

func (e *PrintStmt) isAstNode() {}
func (e *PrintStmt) isStmt() {}

type VarStmt struct {
	Name scanner.Token
	Init Expr
}

func (e *VarStmt) isAstNode() {}
func (e *VarStmt) isStmt() {}

type ExprStmt struct {
	Expr Expr
}

func (e *ExprStmt) isAstNode() {}
func (e *ExprStmt) isStmt() {}
