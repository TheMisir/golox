package main

import (
	"fmt"
)

type Interpereter struct{}

func MakeInterpereter() *Interpereter {
	return &Interpereter{}
}

func (i *Interpereter) evaulate(expr Expr) Any {
	return expr.accept(i)
}

func (i *Interpereter) visitBinaryExpr(expr *BinaryExpr) Any {
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

func (i *Interpereter) visitGroupingExpr(expr *GroupingExpr) Any {
	return i.evaulate(expr.expression)
}

func (i *Interpereter) visitLiteralExpr(expr *LiteralExpr) Any {
	return expr.value
}

func (i *Interpereter) visitUnaryExpr(expr *UnaryExpr) Any {
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

func (i *Interpereter) checkNumberOperand(operator *Token, operand Any) {
	switch operand.(type) {
	case float:
		return
	}
	i.runtimeError(operator, "Operand must be a number.")
}

func (i *Interpereter) checkNumberOperands(operator *Token, left Any, right Any) {
	switch left.(type) {
	case float:
		switch right.(type) {
		case float:
			return
		}
	}
	i.runtimeError(operator, "Operands must be a numbers.")
}

func (i *Interpereter) runtimeError(token *Token, message string) {
	panic(RuntimeError{token: token, message: message})
}

type RuntimeError struct {
	token   *Token
	message string
}

func (e RuntimeError) Error() string {
	if e.token == nil {
		return "Runtime error: " + e.message
	}
	return fmt.Sprintf("Runtime error at line %v: %s", e.token.line, e.message)
}
