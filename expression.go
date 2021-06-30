package main

type BinaryExpr struct {
	left     Expr
	operator *Token
	right    Expr
}

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

func MakeBinaryExpr(left Expr, operator *Token, right Expr) *BinaryExpr {
	return &BinaryExpr{left: left, operator: operator, right: right}
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

func (expr *BinaryExpr) accept(v ExprVisitor) Any {
	return v.visitBinaryExpr(expr)
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
