package parser

import "fmt"

type Visitor interface {
	Visit(node ASTNode) (v Visitor)
}

type PrefixPrinter struct {
	HadError bool
}

func (p *PrefixPrinter) Visit(node ASTNode) (v Visitor) {
	switch n := node.(type) {
	case *FunStmt:
		fmt.Println("function " + n.Name.Lexeme)
	case *IfStmt:

	}
	return v
}

func Walk(v Visitor, node ASTNode) {
	if v := v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *FunStmt:
		Walk(v, &n.Fun)
	case *IfStmt:
		Walk(v, n.Condition)
		Walk(v, n.ThenBranch)
		Walk(v, n.ElseBranch)
	case *WhileStmt:
		Walk(v, n.Condition)
		Walk(v, n.Body)
	case *ExprStmt:
		Walk(v, n.Expr)
	case *VarStmt:
		Walk(v, n.Init)
	case *BinExpr:
		Walk(v, n.Left)
		Walk(v, n.Right)
	case *LogicalExpr:
		Walk(v, n.Left)
		Walk(v, n.Right)
	case *UnExpr:
		Walk(v, n.Right)
	case *LiteralExpr:
		Walk(v, n)
	}
}