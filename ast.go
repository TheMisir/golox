package main

type GetExpr struct {
	object Expr
	name   *Token
}

type ExpressionStmt struct {
	expression Expr
}

type PrintStmt struct {
	expression Expr
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

type ReturnStmt struct {
	keyword *Token
	value   Expr
}

type LiteralExpr struct {
	value interface{}
}

type LogicalExpr struct {
	left     Expr
	operator *Token
	right    Expr
}

type VariableExpr struct {
	name *Token
}

type ThisExpr struct {
	keyword *Token
}

type BlockStmt struct {
	statements []Stmt
}

type ClassStmt struct {
	name    *Token
	methods []*FunctionStmt
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

type FunctionStmt struct {
	name   *Token
	params []*Token
	body   []Stmt
}

type AssignExpr struct {
	name  *Token
	value Expr
}

type CallExpr struct {
	callee    Expr
	paren     *Token
	arguments []Expr
}

type SetExpr struct {
	object Expr
	name   *Token
	value  Expr
}

type VarStmt struct {
	name        *Token
	initializer Expr
}

type BinaryExpr struct {
	left     Expr
	operator *Token
	right    Expr
}

type GroupingExpr struct {
	expression Expr
}

type UnaryExpr struct {
	operator *Token
	right    Expr
}

func MakeGetExpr(object Expr, name *Token) *GetExpr {
	return &GetExpr{object: object, name: name}
}

func MakeExpressionStmt(expression Expr) *ExpressionStmt {
	return &ExpressionStmt{expression: expression}
}

func MakePrintStmt(expression Expr) *PrintStmt {
	return &PrintStmt{expression: expression}
}

func MakeWhileStmt(condition Expr, body Stmt) *WhileStmt {
	return &WhileStmt{condition: condition, body: body}
}

func MakeReturnStmt(keyword *Token, value Expr) *ReturnStmt {
	return &ReturnStmt{keyword: keyword, value: value}
}

func MakeLiteralExpr(value interface{}) *LiteralExpr {
	return &LiteralExpr{value: value}
}

func MakeLogicalExpr(left Expr, operator *Token, right Expr) *LogicalExpr {
	return &LogicalExpr{left: left, operator: operator, right: right}
}

func MakeVariableExpr(name *Token) *VariableExpr {
	return &VariableExpr{name: name}
}

func MakeThisExpr(keyword *Token) *ThisExpr {
	return &ThisExpr{keyword: keyword}
}

func MakeBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{statements: statements}
}

func MakeClassStmt(name *Token, methods []*FunctionStmt) *ClassStmt {
	return &ClassStmt{name: name, methods: methods}
}

func MakeIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) *IfStmt {
	return &IfStmt{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

func MakeFunctionStmt(name *Token, params []*Token, body []Stmt) *FunctionStmt {
	return &FunctionStmt{name: name, params: params, body: body}
}

func MakeAssignExpr(name *Token, value Expr) *AssignExpr {
	return &AssignExpr{name: name, value: value}
}

func MakeCallExpr(callee Expr, paren *Token, arguments []Expr) *CallExpr {
	return &CallExpr{callee: callee, paren: paren, arguments: arguments}
}

func MakeSetExpr(object Expr, name *Token, value Expr) *SetExpr {
	return &SetExpr{object: object, name: name, value: value}
}

func MakeVarStmt(name *Token, initializer Expr) *VarStmt {
	return &VarStmt{name: name, initializer: initializer}
}

func MakeBinaryExpr(left Expr, operator *Token, right Expr) *BinaryExpr {
	return &BinaryExpr{left: left, operator: operator, right: right}
}

func MakeGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{expression: expression}
}

func MakeUnaryExpr(operator *Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator: operator, right: right}
}

func (expr *GetExpr) accept(v ExprVisitor) Any {
	return v.visitGetExpr(expr)
}

func (expr *ExpressionStmt) accept(v StmtVisitor) Any {
	return v.visitExpressionStmt(expr)
}

func (expr *PrintStmt) accept(v StmtVisitor) Any {
	return v.visitPrintStmt(expr)
}

func (expr *WhileStmt) accept(v StmtVisitor) Any {
	return v.visitWhileStmt(expr)
}

func (expr *ReturnStmt) accept(v StmtVisitor) Any {
	return v.visitReturnStmt(expr)
}

func (expr *LiteralExpr) accept(v ExprVisitor) Any {
	return v.visitLiteralExpr(expr)
}

func (expr *LogicalExpr) accept(v ExprVisitor) Any {
	return v.visitLogicalExpr(expr)
}

func (expr *VariableExpr) accept(v ExprVisitor) Any {
	return v.visitVariableExpr(expr)
}

func (expr *ThisExpr) accept(v ExprVisitor) Any {
	return v.visitThisExpr(expr)
}

func (expr *BlockStmt) accept(v StmtVisitor) Any {
	return v.visitBlockStmt(expr)
}

func (expr *ClassStmt) accept(v StmtVisitor) Any {
	return v.visitClassStmt(expr)
}

func (expr *IfStmt) accept(v StmtVisitor) Any {
	return v.visitIfStmt(expr)
}

func (expr *FunctionStmt) accept(v StmtVisitor) Any {
	return v.visitFunctionStmt(expr)
}

func (expr *AssignExpr) accept(v ExprVisitor) Any {
	return v.visitAssignExpr(expr)
}

func (expr *CallExpr) accept(v ExprVisitor) Any {
	return v.visitCallExpr(expr)
}

func (expr *SetExpr) accept(v ExprVisitor) Any {
	return v.visitSetExpr(expr)
}

func (expr *VarStmt) accept(v StmtVisitor) Any {
	return v.visitVarStmt(expr)
}

func (expr *BinaryExpr) accept(v ExprVisitor) Any {
	return v.visitBinaryExpr(expr)
}

func (expr *GroupingExpr) accept(v ExprVisitor) Any {
	return v.visitGroupingExpr(expr)
}

func (expr *UnaryExpr) accept(v ExprVisitor) Any {
	return v.visitUnaryExpr(expr)
}
