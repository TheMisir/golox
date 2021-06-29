package main

type BinaryExpr struct {
  left Expr
  operator *Token
  right Expr
}

type GroupingExpr struct {
  expression Expr
}

type LiteralExpr struct {
  value interface{}
}

type UnaryExpr struct {
  operator *Token
  right Expr
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

