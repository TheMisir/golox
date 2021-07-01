package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	types := map[string]string{
		// Expressions
		"AssignExpr":   "name *Token, value Expr",
		"BinaryExpr":   "left Expr, operator *Token, right Expr",
		"CallExpr":     "callee Expr, paren *Token, arguments []Expr",
		"GetExpr":      "object Expr, name *Token",
		"GroupingExpr": "expression Expr",
		"LiteralExpr":  "value interface{}",
		"LogicalExpr":  "left Expr, operator *Token, right Expr",
		"SetExpr":      "object Expr, name *Token, value Expr",
		"UnaryExpr":    "operator *Token, right Expr",
		"VariableExpr": "name *Token",

		// Statements
		"BlockStmt":      "statements []Stmt",
		"ClassStmt":      "name *Token, methods []*FunctionStmt",
		"ExpressionStmt": "expression Expr",
		"IfStmt":         "condition Expr, thenBranch Stmt, elseBranch Stmt",
		"PrintStmt":      "expression Expr",
		"VarStmt":        "name *Token, initializer Expr",
		"WhileStmt":      "condition Expr, body Stmt",
		"FunctionStmt":   "name *Token, params []*Token, body []Stmt",
		"ReturnStmt":     "keyword *Token, value Expr",
	}

	defs := "package main\n\n"
	impl := ""
	ctor := ""

	for name, params := range types {
		params := strings.Split(params, ",")
		typeName := name[len(name)-4:]

		defs += fmt.Sprintf("type %s struct {\n", name)
		ctor += fmt.Sprintf("func Make%s(", name)

		args := ""

		for i, param := range params {
			param := strings.Trim(param, " ")
			paramName := strings.Split(param, " ")[0]
			defs += fmt.Sprintf("  %s\n", param)
			if i > 0 {
				ctor += ", "
				args += ", "
			}
			ctor += param
			args += fmt.Sprintf("%s: %s", paramName, paramName)
		}
		defs += "}\n\n"
		ctor += fmt.Sprintf(") *%s {\n  return &%s{%s}\n}\n\n", name, name, args)
		impl += fmt.Sprintf("func (expr *%s) accept(v %sVisitor) Any {\n  return v.visit%s(expr)\n}\n\n", name, typeName, name)
	}

	os.Stdout.WriteString(defs)
	os.Stdout.WriteString(ctor)
	os.Stdout.WriteString(impl)
}
