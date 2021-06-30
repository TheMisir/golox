package main

type VarStmt struct {
	name        *Token
	initializer Expr
}

type GroupingExpr struct {
	expression Expr
}

type LogicalExpr struct {
	left     Expr
	operator *Token
	right    Expr
}

type VariableExpr struct {
	name *Token
}

type ExpressionStmt struct {
	expression Expr
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

type PrintStmt struct {
	expression Expr
}

type AssignExpr struct {
	name  *Token
	value Expr
}

type BinaryExpr struct {
	left     Expr
	operator *Token
	right    Expr
}

type LiteralExpr struct {
	value interface{}
}

type UnaryExpr struct {
	operator *Token
	right    Expr
}

type BlockStmt struct {
	statements []Stmt
}

func MakeVarStmt(name *Token, initializer Expr) *VarStmt {
	return &VarStmt{name: name, initializer: initializer}
}

func MakeGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{expression: expression}
}

func MakeLogicalExpr(left Expr, operator *Token, right Expr) *LogicalExpr {
	return &LogicalExpr{left: left, operator: operator, right: right}
}

func MakeVariableExpr(name *Token) *VariableExpr {
	return &VariableExpr{name: name}
}

func MakeExpressionStmt(expression Expr) *ExpressionStmt {
	return &ExpressionStmt{expression: expression}
}

func MakeIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) *IfStmt {
	return &IfStmt{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

func MakePrintStmt(expression Expr) *PrintStmt {
	return &PrintStmt{expression: expression}
}

func MakeAssignExpr(name *Token, value Expr) *AssignExpr {
	return &AssignExpr{name: name, value: value}
}

func MakeBinaryExpr(left Expr, operator *Token, right Expr) *BinaryExpr {
	return &BinaryExpr{left: left, operator: operator, right: right}
}

func MakeLiteralExpr(value interface{}) *LiteralExpr {
	return &LiteralExpr{value: value}
}

func MakeUnaryExpr(operator *Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator: operator, right: right}
}

func MakeBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{statements: statements}
}

func (expr *VarStmt) accept(v StmtVisitor) Any {
	return v.visitVarStmt(expr)
}

func (expr *GroupingExpr) accept(v ExprVisitor) Any {
	return v.visitGroupingExpr(expr)
}

func (expr *LogicalExpr) accept(v ExprVisitor) Any {
	return v.visitLogicalExpr(expr)
}

func (expr *VariableExpr) accept(v ExprVisitor) Any {
	return v.visitVariableExpr(expr)
}

func (expr *ExpressionStmt) accept(v StmtVisitor) Any {
	return v.visitExpressionStmt(expr)
}

func (expr *IfStmt) accept(v StmtVisitor) Any {
	return v.visitIfStmt(expr)
}

func (expr *PrintStmt) accept(v StmtVisitor) Any {
	return v.visitPrintStmt(expr)
}

func (expr *AssignExpr) accept(v ExprVisitor) Any {
	return v.visitAssignExpr(expr)
}

func (expr *BinaryExpr) accept(v ExprVisitor) Any {
	return v.visitBinaryExpr(expr)
}

func (expr *LiteralExpr) accept(v ExprVisitor) Any {
	return v.visitLiteralExpr(expr)
}

func (expr *UnaryExpr) accept(v ExprVisitor) Any {
	return v.visitUnaryExpr(expr)
}

func (expr *BlockStmt) accept(v StmtVisitor) Any {
	return v.visitBlockStmt(expr)
}
