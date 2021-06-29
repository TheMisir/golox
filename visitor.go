package main

import "fmt"

type Any = interface{}

type ExprVisitor interface {
	visitBinaryExpr(expr *BinaryExpr) Any
	visitGroupingExpr(expr *GroupingExpr) Any
	visitLiteralExpr(expr *LiteralExpr) Any
	visitUnaryExpr(expr *UnaryExpr) Any
}

type Expr interface {
	accept(v ExprVisitor) Any
}

type AstPrinter struct{}

func (p *AstPrinter) print(expr Expr) string {
	return expr.accept(p).(string)
}

func (p *AstPrinter) visitBinaryExpr(b *BinaryExpr) Any {
	return fmt.Sprintf("BinaryExpr(%s %s %s)", p.print(b.left), b.operator.tokenType, p.print(b.right))
}

func (p *AstPrinter) visitGroupingExpr(g *GroupingExpr) Any {
	return fmt.Sprintf("GroupingExpr(%s)", p.print(g.expression))
}

func (p *AstPrinter) visitLiteralExpr(l *LiteralExpr) Any {
	return fmt.Sprintf("LiteralExpr(%v)", l.value)
}

func (p *AstPrinter) visitUnaryExpr(u *UnaryExpr) Any {
	return fmt.Sprintf("UnaryExpr(%s %s)", u.operator.tokenType, p.print(u.right))
}
