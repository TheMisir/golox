package main

import (
	"fmt"
	"os"
	"time"
)

type Interpreter struct {
	context     *LoxContext
	environment *Environment
	globals     *Environment
	locals      map[Expr]int
}

func MakeInterpreter(context *LoxContext) *Interpreter {
	globals := MakeEnvironment(context, nil)

	globals.define("clock", MakeLoxCallable(0, func(interpreter *Interpreter, arguments []Any) Any {
		return time.Now().UnixNano() / 1000000
	}))

	return &Interpreter{
		context:     context,
		environment: globals,
		globals:     globals,
		locals:      make(map[Expr]int),
	}
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
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

	distance, ok := i.locals[expr]
	if ok {
		i.environment.assignAt(distance, expr.name.lexme, value)
	} else {
		i.globals.assign(expr.name, value)
	}

	return value
}

func (i *Interpreter) visitVariableExpr(expr *VariableExpr) Any {
	return i.lookUpVariable(expr.name, expr)
}

func (i *Interpreter) lookUpVariable(name *Token, expr Expr) Any {
	distance, ok := i.locals[expr]

	if ok {
		return i.environment.getAt(distance, name.lexme)
	} else {
		return i.globals.get(name)
	}
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

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []Any) Any
	Arity() int
}

type LoxCallableHandler = func(interpreter *Interpreter, arguments []Any) Any

type LoxStaticCallable struct {
	arity   int
	handler LoxCallableHandler
}

func MakeLoxCallable(arity int, handler LoxCallableHandler) *LoxStaticCallable {
	return &LoxStaticCallable{arity: arity, handler: handler}
}

func (c *LoxStaticCallable) Call(interpreter *Interpreter, arguments []Any) Any {
	return c.handler(interpreter, arguments)
}

func (c *LoxStaticCallable) Arity() int {
	return c.arity
}

func (i *Interpreter) visitCallExpr(expr *CallExpr) Any {
	callee := i.evaluate(expr.callee)

	switch val := callee.(type) {
	case LoxCallable:
		if val.Arity() != len(expr.arguments) {
			i.context.runtimeError(expr.paren, "Expected %v arguments byt got %v.", val.Arity(), len(expr.arguments))
		}

		arguments := make([]Any, len(expr.arguments))
		for index, argument := range expr.arguments {
			arguments[index] = i.evaluate(argument)
		}

		return val.Call(i, arguments)
	default:
		i.context.runtimeError(expr.paren, "Can only call functions and classes.")
		return nil
	}
}

type LoxFunction struct {
	declaration *FunctionStmt
	closure     *Environment
}

func MakeLoxFunction(declaration *FunctionStmt, closure *Environment) *LoxFunction {
	return &LoxFunction{declaration: declaration, closure: closure}
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []Any) (result Any) {
	environment := MakeEnvironment(interpreter.context, f.closure)
	for index, param := range f.declaration.params {
		environment.define(param.lexme, arguments[index])
	}

	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case *Return:
				result = r.value
				break
			default:
				panic(r)
			}
		}
	}()

	interpreter.executeBlock(f.declaration.body, environment)
	return nil
}

func (i *Interpreter) visitFunctionStmt(stmt *FunctionStmt) Any {
	function := MakeLoxFunction(stmt, i.environment)
	i.environment.define(stmt.name.lexme, function)
	return nil
}

type Return struct {
	value Any
}

func MakeReturn(value Any) *Return {
	return &Return{value: value}
}

func (r *Return) Error() string {
	return "return statement"
}

func (i *Interpreter) visitReturnStmt(stmt *ReturnStmt) Any {
	var value Any = nil
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	}

	panic(MakeReturn(value))
}

func (i *Interpreter) visitClassStmt(stmt *ClassStmt) Any {
	i.environment.define(stmt.name.lexme, nil)
	klass := MakeLoxClass(stmt.name.lexme)
	i.environment.assign(stmt.name, klass)
	return nil
}

type LoxClass struct {
	name string
}

func MakeLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}

func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) Arity() int {
	return 0
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []Any) Any {
	instance := MakeLoxInstance(c)
	return instance
}

type LoxInstance struct {
	klass *LoxClass
}

func MakeLoxInstance(klass *LoxClass) *LoxInstance {
	return &LoxInstance{klass: klass}
}

func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.klass.name)
}
