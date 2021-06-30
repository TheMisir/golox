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
		return left.(float64) - right.(float64)

	case SLASH:
		return left.(float64) / right.(float64)

	case STAR:
		return left.(float64) * right.(float64)

	case PLUS:
		switch leftVal := left.(type) {
		case float64:
			switch rightVal := right.(type) {
			case float64:
				return leftVal + rightVal
			}
		case string:
			switch rightVal := right.(type) {
			case string:
				return leftVal + rightVal
			case float64:
				return fmt.Sprintf("%s%v", leftVal, rightVal)
			}
		}
		break

	case GREATER:
		return left.(float64) > right.(float64)

	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)

	case LESS:
		return left.(float64) < right.(float64)

	case LESS_EQUAL:
		return left.(float64) <= right.(float64)

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
		return -right.(float64)
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
