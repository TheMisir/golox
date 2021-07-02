package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

var treePrinter = &AstPrinter{}

func (p *AstPrinter) print(node Any) string {
	if node == nil {
		return "<nil>"
	}

	switch node := node.(type) {
	case Stmt:
		return node.accept(p).(string)
	case Expr:
		return node.accept(p).(string)
	case []Stmt:
		statements := make([]string, len(node))
		for index, statement := range node {
			statements[index] = p.print(statement)
		}
		return strings.Join(statements, "; \n")
	case []Expr:
		expressions := make([]string, len(node))
		for index, expression := range node {
			expressions[index] = p.print(expression)
		}
		return strings.Join(expressions, ", ")
	case []*Token:
		tokens := make([]string, len(node))
		for index, token := range node {
			tokens[index] = p.print(token)
		}
		return strings.Join(tokens, ", ")
	case []Token:
		tokens := make([]string, len(node))
		for index, token := range node {
			tokens[index] = p.print(token)
		}
		return strings.Join(tokens, ", ")
	case *Token:
		return p.printToken(node)
	case Token:
		return p.printToken(&node)
	default:
		return fmt.Sprintf("<%T>", node)
	}
}

func (p *AstPrinter) printToken(token *Token) string {
	if token == nil {
		return "<nil>"
	}

	switch token.tokenType {
	case IDENTIFIER:
		return token.lexme
	case STRING:
		return fmt.Sprintf("Literal(%v)", token.literal)
	case NUMBER:
		return fmt.Sprintf("Literal(%v)", token.literal)
	default:
		return string(token.tokenType)
	}
}

func (p *AstPrinter) visitBinaryExpr(expr *BinaryExpr) Any {
	return fmt.Sprintf("Binary(%s %s %s)", p.print(expr.left), expr.operator.tokenType, p.print(expr.right))
}

func (p *AstPrinter) visitGroupingExpr(expr *GroupingExpr) Any {
	return fmt.Sprintf("(%s)", p.print(expr.expression))
}

func (p *AstPrinter) visitLiteralExpr(expr *LiteralExpr) Any {
	return fmt.Sprintf("Literal(%v)", expr.value)
}

func (p *AstPrinter) visitUnaryExpr(expr *UnaryExpr) Any {
	return fmt.Sprintf("Unary(%s %s)", p.print(expr.operator), p.print(expr.right))
}

func (p *AstPrinter) visitVariableExpr(expr *VariableExpr) Any {
	return fmt.Sprintf("Variable(%s)", p.print(expr.name))
}

func (p *AstPrinter) visitPrintStmt(stmt *PrintStmt) Any {
	return fmt.Sprintf("Print(%s);", p.print(stmt.expression))
}

func (p *AstPrinter) visitExpressionStmt(stmt *ExpressionStmt) Any {
	return fmt.Sprintf("Expression(%s)", p.print(stmt.expression))
}

func (p *AstPrinter) visitVarStmt(stmt *VarStmt) Any {
	return fmt.Sprintf("Var(%s = %s)", p.print(stmt.name), p.print(stmt.initializer))
}

func (p *AstPrinter) visitAssignExpr(expr *AssignExpr) Any {
	return fmt.Sprintf("Assign(%s = %s)", p.print(expr.name), p.print(expr.value))
}

func (p *AstPrinter) visitBlockStmt(stmt *BlockStmt) Any {
	return fmt.Sprintf("Block(%s)", p.print(stmt.statements))
}

func (p *AstPrinter) visitIfStmt(stmt *IfStmt) Any {
	return fmt.Sprintf("If(%s) {%s} Else {%s}", p.print(stmt.condition), p.print(stmt.thenBranch), p.print(stmt.elseBranch))
}

func (p *AstPrinter) visitLogicalExpr(expr *LogicalExpr) Any {
	return fmt.Sprintf("Logical(%s %s %s)", p.print(expr.left), p.print(expr.operator), p.print(expr.right))
}

func (p *AstPrinter) visitWhileStmt(stmt *WhileStmt) Any {
	return fmt.Sprintf("While(%s) {%s}", p.print(stmt.condition), p.print(stmt.body))
}

func (p *AstPrinter) visitCallExpr(expr *CallExpr) Any {
	arguments := make([]string, len(expr.arguments)+1)
	arguments[0] = p.print(expr.callee)

	for index, arg := range expr.arguments {
		arguments[index+1] = p.print(arg)
	}

	return fmt.Sprintf("Call(%s (%s))", p.print(expr.callee), p.print(expr.arguments))
}

func (p *AstPrinter) visitFunctionExpr(expr *FunctionExpr) Any {
	return fmt.Sprintf("Function(%s (%s)) {%s}", p.print(expr.name), p.print(expr.params), p.print(expr.body))
}

func (p *AstPrinter) visitReturnStmt(stmt *ReturnStmt) Any {
	return fmt.Sprintf("Return(%s)", p.print(stmt.value))
}

func (p *AstPrinter) visitClassStmt(stmt *ClassStmt) Any {
	return fmt.Sprintf("Class(%s < %s) {%s}", p.print(stmt.name), p.print(stmt.superclass), p.print(stmt.methods))
}

func (p *AstPrinter) visitGetExpr(expr *GetExpr) Any {
	return fmt.Sprintf("Get(%s.%s)", p.print(expr.object), p.print(expr.name))
}

func (p *AstPrinter) visitSetExpr(expr *SetExpr) Any {
	return fmt.Sprintf("Set(%s.%s = %s)", p.print(expr.object), p.print(expr.name), p.print(expr.value))
}

func (p *AstPrinter) visitThisExpr(expr *ThisExpr) Any {
	return "This"
}

func (p *AstPrinter) visitSuperExpr(expr *SuperExpr) Any {
	return fmt.Sprintf("Super(%s)", expr.method.lexme)
}

func (p *AstPrinter) visitForStmt(stmt *ForStmt) Any {
	return fmt.Sprintf("For(%s; %s; %s) {%s}", p.print(stmt.initializer), p.print(stmt.condition), p.print(stmt.increment), p.print(stmt.body))
}

func (p *AstPrinter) visitContinueStmt(stmt *ContinueStmt) Any {
	return "Continue"
}

func (p *AstPrinter) visitBreakStmt(stmt *BreakStmt) Any {
	return "Break"
}
