package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

var treePrinter = &AstPrinter{}

func (p *AstPrinter) printExpr(expr Expr) string {
	return expr.accept(p).(string)
}

func (p *AstPrinter) printStmt(stmt Stmt) string {
	return stmt.accept(p).(string)
}

func (p *AstPrinter) visitBinaryExpr(expr *BinaryExpr) Any {
	return fmt.Sprintf("BinaryExpr(%s %s %s)", p.printExpr(expr.left), expr.operator.tokenType, p.printExpr(expr.right))
}

func (p *AstPrinter) visitGroupingExpr(expr *GroupingExpr) Any {
	return fmt.Sprintf("GroupingExpr(%s)", p.printExpr(expr.expression))
}

func (p *AstPrinter) visitLiteralExpr(expr *LiteralExpr) Any {
	return fmt.Sprintf("LiteralExpr(%v)", expr.value)
}

func (p *AstPrinter) visitUnaryExpr(expr *UnaryExpr) Any {
	return fmt.Sprintf("UnaryExpr(%s %s)", expr.operator.tokenType, p.printExpr(expr.right))
}

func (p *AstPrinter) visitVariableExpr(expr *VariableExpr) Any {
	return fmt.Sprintf("VariableExpr(%s)", expr.name.lexme)
}

func (p *AstPrinter) visitPrintStmt(stmt *PrintStmt) Any {
	return fmt.Sprintf("PrintStmt(%s);", p.printExpr(stmt.expression))
}

func (p *AstPrinter) visitExpressionStmt(stmt *ExpressionStmt) Any {
	return fmt.Sprintf("ExpressionStmt(%s);", p.printExpr(stmt.expression))
}

func (p *AstPrinter) visitVarStmt(stmt *VarStmt) Any {
	return fmt.Sprintf("VarStmt(%s %s)", stmt.name.lexme, p.printExpr(stmt.initializer))
}

func (p *AstPrinter) visitAssignExpr(expr *AssignExpr) Any {
	return fmt.Sprintf("AssignExpr(%s = %s)", expr.name.lexme, p.printExpr(expr.value))
}

func (p *AstPrinter) visitBlockStmt(stmt *BlockStmt) Any {
	result := "BlockStmt{"
	for index, statement := range stmt.statements {
		if index > 0 {
			result += " "
		}
		result += p.printStmt(statement)
	}
	return result + "}"
}

func (p *AstPrinter) visitIfStmt(stmt *IfStmt) Any {
	if stmt.elseBranch == nil {
		return fmt.Sprintf("IfStmt(%s {%s})", p.printExpr(stmt.condition), p.printStmt(stmt.thenBranch))
	}

	return fmt.Sprintf("IfStmt(%s {%s} else {%s})", p.printExpr(stmt.condition), p.printStmt(stmt.thenBranch), p.printStmt(stmt.elseBranch))
}

func (p *AstPrinter) visitLogicalExpr(expr *LogicalExpr) Any {
	return fmt.Sprintf("LogicalExpr(%s %s %s)", p.printExpr(expr.left), expr.operator.tokenType, p.printExpr(expr.right))
}

func (p *AstPrinter) visitWhileStmt(stmt *WhileStmt) Any {
	return fmt.Sprintf("WhileStmt(%s {%s})", p.printExpr(stmt.condition), p.printStmt(stmt.body))
}

func (p *AstPrinter) visitCallExpr(expr *CallExpr) Any {
	arguments := make([]string, len(expr.arguments)+1)
	arguments[0] = p.printExpr(expr.callee)

	for index, arg := range expr.arguments {
		arguments[index+1] = p.printExpr(arg)
	}

	return fmt.Sprintf("CallExpr(%s)", strings.Join(arguments, ", "))
}

func (p *AstPrinter) visitFunctionStmt(stmt *FunctionStmt) Any {
	params := make([]string, len(stmt.params))
	for index, param := range stmt.params {
		params[index] = param.String()
	}

	body := make([]string, len(stmt.body))
	for index, element := range stmt.body {
		body[index] = p.printStmt(element)
	}

	return fmt.Sprintf("FunctionStmt(%s(%s) {%s})", stmt.name.lexme, strings.Join(params, ", "), strings.Join(body, "; "))
}

func (p *AstPrinter) visitReturnStmt(stmt *ReturnStmt) Any {
	return fmt.Sprintf("ReturnStmt(%s)", p.printExpr(stmt.value))
}

func (p *AstPrinter) visitClassStmt(stmt *ClassStmt) Any {
	methods := make([]string, len(stmt.methods))
	for index, method := range stmt.methods {
		methods[index] = p.printStmt(method)
	}

	return fmt.Sprintf("ClassStmt(%s) {%s}", stmt.name.lexme, strings.Join(methods, ", "))
}

func (p *AstPrinter) visitGetExpr(expr *GetExpr) Any {
	return fmt.Sprintf("GetExpr(%s.%s)", p.printExpr(expr.object), expr.name.lexme)
}

func (p *AstPrinter) visitSetExpr(expr *SetExpr) Any {
	return fmt.Sprintf("SetExpr(%s.%s = %s)", p.printExpr(expr.object), expr.name.lexme, p.printExpr(expr.value))
}

func (p *AstPrinter) visitThisExpr(expr *ThisExpr) Any {
	return "ThisExpr"
}
