package main

type GroupingExpr struct {
	expression Expr
}

type LiteralExpr struct {
	value interface{}
}

type UnaryExpr struct {
	operator *Token
	right    Expr
}

type ExpressionStmt struct {
	expression Expr
}

type PrintStmt struct {
	expression Expr
}

type BinaryExpr struct {
	left     Expr
	operator *Token
	right    Expr
}

func MakeGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{expression: expression}
}

func MakeLiteralExpr(value interface{}) *LiteralExpr {
	return &LiteralExpr{value: value}
}

func MakeUnaryExpr(operator *Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator: operator, right: right}
}

func MakeExpressionStmt(expression Expr) *ExpressionStmt {
	return &ExpressionStmt{expression: expression}
}

func MakePrintStmt(expression Expr) *PrintStmt {
	return &PrintStmt{expression: expression}
}

func MakeBinaryExpr(left Expr, operator *Token, right Expr) *BinaryExpr {
	return &BinaryExpr{left: left, operator: operator, right: right}
}

func (expr *GroupingExpr) accept(v ExprVisitor) Any {
	return v.visitGroupingExpr(expr)
}

func (expr *LiteralExpr) accept(v ExprVisitor) Any {
	return v.visitLiteralExpr(expr)
}

func (expr *UnaryExpr) accept(v ExprVisitor) Any {
	return v.visitUnaryExpr(expr)
}

func (expr *ExpressionStmt) accept(v StmtVisitor) Any {
	return v.visitExpressionStmt(expr)
}

func (expr *PrintStmt) accept(v StmtVisitor) Any {
	return v.visitPrintStmt(expr)
}

func (expr *BinaryExpr) accept(v ExprVisitor) Any {
	return v.visitBinaryExpr(expr)
}
