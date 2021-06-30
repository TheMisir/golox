package main

import "fmt"

type AstPrinter struct{}

func MakeAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

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
	return fmt.Sprintf("AssignExpr(%s %s)", expr.name.lexme, p.printExpr(expr.value))
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
