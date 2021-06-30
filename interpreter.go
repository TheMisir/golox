package main

import (
	"fmt"
	"os"
)

type Interpreter struct {
	context     *LoxContext
	environment *Environment
}

func MakeInterpreter(context *LoxContext) *Interpreter {
	return &Interpreter{
		context:     context,
		environment: MakeEnvironment(context, nil),
	}
}

func (i *Interpreter) evaluate(expr Expr) Any {
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
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

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
		i.context.runtimeError(expr.operator, "Operands must be two numbers or two strings.")
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
	return i.evaluate(expr.expression)
}

func (i *Interpreter) visitLiteralExpr(expr *LiteralExpr) Any {
	return expr.value
}

func (i *Interpreter) visitUnaryExpr(expr *UnaryExpr) Any {
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		i.checkNumberOperand(expr.operator, right)
		return -right.(float)
	case BANG:
		return !isTruthy(right)
	}

	return nil
}

func (i *Interpreter) visitAssignExpr(expr *AssignExpr) Any {
	value := i.evaluate(expr.value)
	i.environment.assign(expr.name, value)
	return value
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
	i.context.runtimeError(operator, "Operand must be a number.")
}

func (i *Interpreter) checkNumberOperands(operator *Token, left Any, right Any) {
	switch left.(type) {
	case float:
		switch right.(type) {
		case float:
			return
		}
	}
	i.context.runtimeError(operator, "Operands must be a numbers.")
}

func (i *Interpreter) visitPrintStmt(stmt *PrintStmt) Any {
	value := i.evaluate(stmt.expression)
	fmt.Fprintf(os.Stdout, "%v\n", value)
	return nil
}

func (i *Interpreter) visitExpressionStmt(stmt *ExpressionStmt) Any {
	i.evaluate(stmt.expression)
	return nil
}

func (i *Interpreter) visitVarStmt(stmt *VarStmt) Any {
	var value Any = nil
	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}

	i.environment.define(stmt.name.lexme, value)
	return nil
}

func (i *Interpreter) visitBlockStmt(stmt *BlockStmt) Any {
	i.executeBlock(stmt.statements, MakeEnvironment(i.context, i.environment))
	return nil
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()
	i.environment = environment

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) visitIfStmt(stmt *IfStmt) Any {
	if isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}
	return nil
}

func (i *Interpreter) visitLogicalExpr(expr *LogicalExpr) Any {
	left := i.evaluate(expr.left)

	if expr.operator.tokenType == OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.right)
}

func (i *Interpreter) visitWhileStmt(stmt *WhileStmt) Any {
	for isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
	}
	return nil
}
