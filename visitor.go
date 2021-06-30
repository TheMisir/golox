package main

type Any = interface{}
type float = float64

type Expr interface {
	accept(v ExprVisitor) Any
}
type ExprVisitor interface {
	visitBinaryExpr(expr *BinaryExpr) Any
	visitGroupingExpr(expr *GroupingExpr) Any
	visitLiteralExpr(expr *LiteralExpr) Any
	visitUnaryExpr(expr *UnaryExpr) Any
}

type Stmt interface {
	accept(v StmtVisitor) Any
}

type StmtVisitor interface {
	visitExpressionStmt(stmt *ExpressionStmt) Any
	visitPrintStmt(stmt *PrintStmt) Any
}
