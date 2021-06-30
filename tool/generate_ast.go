package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	types := map[string]string{
		"AssignExpr":     "name *Token, value Expr",
		"BinaryExpr":     "left Expr, operator *Token, right Expr",
		"GroupingExpr":   "expression Expr",
		"LiteralExpr":    "value interface{}",
		"UnaryExpr":      "operator *Token, right Expr",
		"VariableExpr":   "name *Token",
		"BlockStmt":      "statements []Stmt",
		"ExpressionStmt": "expression Expr",
		"PrintStmt":      "expression Expr",
		"VarStmt":        "name *Token, initializer Expr",
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
