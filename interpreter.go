package main

import (
	"fmt"
	"os"
)

type Interpreter struct {
	environment *Environment
}

func MakeInterpreter() *Interpreter {
	return &Interpreter{
		environment: MakeEnvironment(),
	}
}

type RuntimeError struct {
	token   *Token
	message string
}

func MakeRuntimeError(token *Token, message string, a ...interface{}) RuntimeError {
	return RuntimeError{token: token, message: fmt.Sprintf(message, a...)}
}

func (e RuntimeError) Error() string {
	if e.token == nil {
		return "Runtime error: " + e.message
	}
	return fmt.Sprintf("Runtime error at line %v: %s", e.token.line, e.message)
}

func (i *Interpreter) runtimeError(token *Token, message string) {
	panic(MakeRuntimeError(token, message))
}

func (i *Interpreter) evaulate(expr Expr) Any {
	return expr.accept(i)
}

func (i *Interpreter) execute(expr Stmt) Any {
	return expr.accept(i)
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) visitBinaryExpr(expr *BinaryExpr) Any {
	left := i.evaulate(expr.left)
	right := i.evaulate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float) - right.(float)

	case SLASH:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float) / right.(float)

	case STAR:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float) * right.(float)

	case PLUS:
		switch leftVal := left.(type) {
		case float:
			switch rightVal := right.(type) {
			case float:
				return leftVal + rightVal
			case string:
				return fmt.Sprintf("%v%s", leftVal, rightVal)
			}
		case string:
			switch rightVal := right.(type) {
			case string:
				return leftVal + rightVal
			case float:
				return fmt.Sprintf("%s%v", leftVal, rightVal)
			}
		}
		i.runtimeError(expr.operator, "Operands must be two numbers or two strings.")
		break

	case GREATER:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float) > right.(float)

	case GREATER_EQUAL:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float) >= right.(float)

	case LESS:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float) < right.(float)

	case LESS_EQUAL:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float) <= right.(float)

	case EQUAL_EQUAL:
		return isEqual(left, right)

	case BANG_EQUAL:
		return !isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) visitGroupingExpr(expr *GroupingExpr) Any {
	return i.evaulate(expr.expression)
}

func (i *Interpreter) visitLiteralExpr(expr *LiteralExpr) Any {
	return expr.value
}

func (i *Interpreter) visitUnaryExpr(expr *UnaryExpr) Any {
	right := i.evaulate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		i.checkNumberOperand(expr.operator, right)
		return -right.(float)
	case BANG:
		return !isTruthy(right)
	}

	return nil
}

func (i *Interpreter) visitVariableExpr(expr *VariableExpr) Any {
	return i.environment.get(expr.name)
}

func isTruthy(value Any) bool {
	if value == nil {
		return false
	}
	if value == false {
		return false
	}
	return true
}

func isEqual(a Any, b Any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) checkNumberOperand(operator *Token, operand Any) {
	switch operand.(type) {
	case float:
		return
	}
	i.runtimeError(operator, "Operand must be a number.")
}

func (i *Interpreter) checkNumberOperands(operator *Token, left Any, right Any) {
	switch left.(type) {
	case float:
		switch right.(type) {
		case float:
			return
		}
	}
	i.runtimeError(operator, "Operands must be a numbers.")
}

func (i *Interpreter) visitPrintStmt(stmt *PrintStmt) Any {
	value := i.evaulate(stmt.expression)
	fmt.Fprintf(os.Stdout, "%v\n", value)
	return nil
}

func (i *Interpreter) visitExpressionStmt(stmt *ExpressionStmt) Any {
	i.evaulate(stmt.expression)
	return nil
}

func (i *Interpreter) visitVarStmt(stmt *VarStmt) Any {
	var value Any = nil
	if stmt.initializer != nil {
		value = i.evaulate(stmt.initializer)
	}

	i.environment.define(stmt.name.lexme, value)
	return nil
}
