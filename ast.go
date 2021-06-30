package main

type AssignExpr struct {
	name  *Token
	value Expr
}

type GroupingExpr struct {
	expression Expr
}

type LiteralExpr struct {
	value interface{}
}

type BlockStmt struct {
	statements []Stmt
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

type UnaryExpr struct {
	operator *Token
	right    Expr
}

type VariableExpr struct {
	name *Token
}

type VarStmt struct {
	name        *Token
	initializer Expr
}

func MakeAssignExpr(name *Token, value Expr) *AssignExpr {
	return &AssignExpr{name: name, value: value}
}

func MakeGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{expression: expression}
}

func MakeLiteralExpr(value interface{}) *LiteralExpr {
	return &LiteralExpr{value: value}
}

func MakeBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{statements: statements}
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

func MakeUnaryExpr(operator *Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator: operator, right: right}
}

func MakeVariableExpr(name *Token) *VariableExpr {
	return &VariableExpr{name: name}
}

func MakeVarStmt(name *Token, initializer Expr) *VarStmt {
	return &VarStmt{name: name, initializer: initializer}
}

func (expr *AssignExpr) accept(v ExprVisitor) Any {
	return v.visitAssignExpr(expr)
}

func (expr *GroupingExpr) accept(v ExprVisitor) Any {
	return v.visitGroupingExpr(expr)
}

func (expr *LiteralExpr) accept(v ExprVisitor) Any {
	return v.visitLiteralExpr(expr)
}

func (expr *BlockStmt) accept(v StmtVisitor) Any {
	return v.visitBlockStmt(expr)
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

func (expr *UnaryExpr) accept(v ExprVisitor) Any {
	return v.visitUnaryExpr(expr)
}

func (expr *VariableExpr) accept(v ExprVisitor) Any {
	return v.visitVariableExpr(expr)
}

func (expr *VarStmt) accept(v StmtVisitor) Any {
	return v.visitVarStmt(expr)
}
