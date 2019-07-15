package parser

import "lox/scanner"

type Expr interface {
	ASTNode
	isExpr()
}

type TernExpr struct {
	Left Expr
	Op scanner.Token
	Right Expr
	Right2 Expr
}

func (e *TernExpr) isAstNode() {}
func (e *TernExpr) isExpr() {}

type BinExpr struct {
	Left Expr
	Op scanner.Token
	Right Expr
}

func (e *BinExpr) isAstNode() {}
func (e *BinExpr) isExpr() {}

type UnExpr struct {
	Op scanner.Token
	Right Expr
}

func (e *UnExpr) isAstNode() {}
func (e *UnExpr) isExpr() {}

type LiteralExpr struct {
	Literal interface{}
}

func (e *LiteralExpr) isAstNode() {}
func (e *LiteralExpr) isExpr() {}

type AssignExpr struct {
	Name scanner.Token
	Right Expr
}

func (e *AssignExpr) isAstNode() {}
func (e *AssignExpr) isExpr() {}

type LogicalExpr struct {
	Left Expr
	Op scanner.Token
	Right Expr
}

func (e *LogicalExpr) isAstNode() {}
func (e *LogicalExpr) isExpr() {}

type CallExpr struct {
	Calle Expr
	Args []Expr
}

func (e *CallExpr) isAstNode() {}
func (e *CallExpr) isExpr() {}

type FunExpr struct {
	Body []Stmt
}

func (e *FunExpr) isAstNode() {}
func (e *FunExpr) isExpr() {}

type GroupExpr struct {
	Expr Expr
}

func (e *GroupExpr) isAstNode() {}
func (e *GroupExpr) isExpr() {}

type VarExpr struct {
	Name scanner.Token
}

func (e *VarExpr) isAstNode() {}
func (e *VarExpr) isExpr() {}
