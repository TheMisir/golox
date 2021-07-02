package main

import (
	"fmt"
	"os"
)

type LoopBodyResult = string

const (
	LOOP_CONTINUE = "CONTINUE"
	LOOP_BREAK    = "BREAK"
)

type Interpreter struct {
	context     *LoxContext
	environment *Environment
	globals     *Environment
	locals      map[Expr]int
	includes    map[Stmt]*Source
}

func MakeInterpreter(context *LoxContext) *Interpreter {
	globals := MakeEnvironment(context, nil)

	InitializeStdLib(globals)

	return &Interpreter{
		context:     context,
		environment: globals,
		globals:     globals,
		locals:      make(map[Expr]int),
		includes:    make(map[Stmt]*Source),
	}
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) include(stmt Stmt, source *Source) {
	i.includes[stmt] = source
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
	i.executeBlock(stmt.statements, i.environment.extend())
	return nil
}

func (i *Interpreter) visitIncludeStmt(stmt *IncludeStmt) Any {
	source := i.includes[stmt]
	i.executeBlock(source.Body, i.environment)
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
		if i.loopBody(stmt.body) == LOOP_BREAK {
			break
		}
	}
	return nil
}

func (i *Interpreter) visitForStmt(stmt *ForStmt) Any {
	var initializer = func() {
		if stmt.initializer != nil {
			i.execute(stmt.initializer)
		}
	}

	var condition = func() bool {
		if stmt.condition != nil {
			return isTruthy(i.evaluate(stmt.condition))
		}
		return true
	}

	var increment = func() {
		if stmt.increment != nil {
			i.evaluate(stmt.increment)
		}
	}

	for initializer(); condition(); increment() {
		if i.loopBody(stmt.body) == LOOP_BREAK {
			break
		}
	}
	return nil
}

type LoopJump struct {
	jumpType LoopBodyResult
}

func (j LoopJump) String() string {
	return j.jumpType
}

func (i *Interpreter) loopBody(body Stmt) (result LoopBodyResult) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case LoopJump:
				result = r.jumpType
				break
			default:
				panic(r)
			}
		}
	}()

	i.execute(body)
	return LOOP_CONTINUE
}

func (i *Interpreter) visitBreakStmt(stmt *BreakStmt) Any {
	panic(LoopJump{LOOP_BREAK})
}

func (i *Interpreter) visitContinueStmt(stmt *ContinueStmt) Any {
	panic(LoopJump{LOOP_CONTINUE})
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
			i.context.runtimeError(expr.paren, "Expected %v arguments but got %v.", val.Arity(), len(expr.arguments))
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
	declaration   *FunctionExpr
	closure       *Environment
	isInitializer bool
}

func MakeLoxFunction(declaration *FunctionExpr, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	environment := f.closure.extend()
	environment.define("this", instance)
	return MakeLoxFunction(f.declaration, environment, f.isInitializer)
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []Any) (result Any) {
	environment := f.closure.extend()
	for index, param := range f.declaration.params {
		environment.define(param.lexme, arguments[index])
	}

	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case *Return:
				if f.isInitializer {
					result = f.closure.getAt(0, "this")
				} else {
					result = r.value
				}
				break
			default:
				panic(r)
			}
		}
	}()

	interpreter.executeBlock(f.declaration.body, environment)

	if f.isInitializer {
		return f.closure.getAt(0, "this")
	}

	return nil
}

func (i *Interpreter) visitFunctionExpr(expr *FunctionExpr) Any {
	function := MakeLoxFunction(expr, i.environment, false)
	if expr.name != nil {
		i.environment.define(expr.name.lexme, function)
	}

	return function
}

type Return struct {
	value Any
}

func (i *Interpreter) visitReturnStmt(stmt *ReturnStmt) Any {
	var value Any = nil
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	}

	panic(&Return{value})
}

func (i *Interpreter) visitClassStmt(stmt *ClassStmt) Any {
	var superclass *LoxClass = nil
	if stmt.superclass != nil {
		switch value := i.evaluate(stmt.superclass).(type) {
		case *LoxClass:
			superclass = value
			break
		default:
			i.context.runtimeError(stmt.superclass.name, "Superclass must be a class.")
		}
	}

	i.environment.define(stmt.name.lexme, nil)

	if superclass != nil {
		i.environment = i.environment.extend()
		i.environment.define("super", superclass)
	}

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.methods {
		function := MakeLoxFunction(method, i.environment, method.name.lexme == "init")
		methods[method.name.lexme] = function
	}

	klass := MakeLoxClass(i.context, stmt.name.lexme, superclass, methods)

	if superclass != nil {
		i.environment = i.environment.enclosing
	}

	i.environment.assign(stmt.name, klass)
	return nil
}

type LoxClass struct {
	context    *LoxContext
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

func MakeLoxClass(context *LoxContext, name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		context:    context,
		name:       name,
		superclass: superclass,
		methods:    methods,
	}
}

func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) Arity() int {
	if initializer := c.findMethod("init"); initializer != nil {
		return initializer.Arity()
	}

	return 0
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []Any) Any {
	instance := MakeLoxInstance(c)
	if initializer := c.findMethod("init"); initializer != nil {
		initializer.bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (c *LoxClass) findMethod(name string) *LoxFunction {
	if method, ok := c.methods[name]; ok {
		return method
	}

	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}

	return nil
}

type LoxInstance struct {
	klass  *LoxClass
	fields map[string]Any
}

func MakeLoxInstance(klass *LoxClass) *LoxInstance {
	return &LoxInstance{
		klass:  klass,
		fields: make(map[string]Any),
	}
}

func (i *LoxInstance) get(name *Token) Any {
	if value, ok := i.fields[name.lexme]; ok {
		return value
	}

	if method := i.klass.findMethod(name.lexme); method != nil {
		return method.bind(i)
	}

	i.klass.context.runtimeError(name, "Undefined property '%s'.", name.lexme)
	return nil
}

func (i *LoxInstance) set(name *Token, value Any) {
	i.fields[name.lexme] = value
}

func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.klass.name)
}

func (i *Interpreter) visitGetExpr(expr *GetExpr) Any {
	object := i.evaluate(expr.object)
	switch object := object.(type) {
	case *LoxInstance:
		return object.get(expr.name)
	}

	i.context.runtimeError(expr.name, "Only instances have properties.")
	return nil
}

func (i *Interpreter) visitSetExpr(expr *SetExpr) Any {
	object := i.evaluate(expr.object)

	switch object := object.(type) {
	case *LoxInstance:
		value := i.evaluate(expr.value)
		object.set(expr.name, value)
		return nil
	}

	i.context.runtimeError(expr.name, "Only instances have fields.")
	return nil
}

func (i *Interpreter) visitThisExpr(expr *ThisExpr) Any {
	return i.lookUpVariable(expr.keyword, expr)
}

func (i *Interpreter) visitSuperExpr(expr *SuperExpr) Any {
	distance := i.locals[expr]
	superclass := i.environment.getAt(distance, "super").(*LoxClass)

	object := i.environment.getAt(distance-1, "this").(*LoxInstance)

	method := superclass.findMethod(expr.method.lexme)
	if method == nil {
		i.context.runtimeError(expr.method, "Undefined property '%s'.", expr.method.lexme)
	}

	return method.bind(object)
}
