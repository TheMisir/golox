package main

type Any = interface{}
type float = float64

type Expr interface {
	accept(v ExprVisitor) Any
}
type ExprVisitor interface {
	visitBinaryExpr(expr *BinaryExpr) Any
	visitCallExpr(expr *CallExpr) Any
	visitGetExpr(expr *GetExpr) Any
	visitSetExpr(expr *SetExpr) Any
	visitGroupingExpr(expr *GroupingExpr) Any
	visitLiteralExpr(expr *LiteralExpr) Any
	visitUnaryExpr(expr *UnaryExpr) Any
	visitVariableExpr(expr *VariableExpr) Any
	visitAssignExpr(expr *AssignExpr) Any
	visitThisExpr(expr *ThisExpr) Any
	visitSuperExpr(expr *SuperExpr) Any
	visitLogicalExpr(expr *LogicalExpr) Any
	visitFunctionExpr(expr *FunctionExpr) Any
}

type Stmt interface {
	accept(v StmtVisitor) Any
}

type StmtVisitor interface {
	visitExpressionStmt(stmt *ExpressionStmt) Any
	visitPrintStmt(stmt *PrintStmt) Any
	visitVarStmt(stmt *VarStmt) Any
	visitBlockStmt(stmt *BlockStmt) Any
	visitIfStmt(stmt *IfStmt) Any
	visitWhileStmt(stmt *WhileStmt) Any
	visitReturnStmt(stmt *ReturnStmt) Any
	visitClassStmt(stmt *ClassStmt) Any
}
