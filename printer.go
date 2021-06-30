package main

import "fmt"

type AstPrinter struct{}

func MakeAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (p *AstPrinter) print(expr Expr) string {
	return expr.accept(p).(string)
}

func (p *AstPrinter) visitBinaryExpr(expr *BinaryExpr) Any {
	return fmt.Sprintf("BinaryExpr(%s %s %s)", p.print(expr.left), expr.operator.tokenType, p.print(expr.right))
}

func (p *AstPrinter) visitGroupingExpr(expr *GroupingExpr) Any {
	return fmt.Sprintf("GroupingExpr(%s)", p.print(expr.expression))
}

func (p *AstPrinter) visitLiteralExpr(expr *LiteralExpr) Any {
	return fmt.Sprintf("LiteralExpr(%v)", expr.value)
}

func (p *AstPrinter) visitUnaryExpr(expr *UnaryExpr) Any {
	return fmt.Sprintf("UnaryExpr(%s %s)", expr.operator.tokenType, p.print(expr.right))
}

func (p *AstPrinter) visitPrintStmt(stmt *PrintStmt) Any {
	return fmt.Sprintf("PrintStmt(%s)")
}

func (p *AstPrinter) visitExpressionStmt(stmt *ExpressionStmt) Any {
	return stmt.accept(p)
}
