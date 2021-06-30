package main

type Any = interface{}
type float = float64

type ExprVisitor interface {
	visitBinaryExpr(expr *BinaryExpr) Any
	visitGroupingExpr(expr *GroupingExpr) Any
	visitLiteralExpr(expr *LiteralExpr) Any
	visitUnaryExpr(expr *UnaryExpr) Any
}

type Expr interface {
	accept(v ExprVisitor) Any
}
