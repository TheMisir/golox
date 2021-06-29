package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	types := []string{
		"BinaryExpr   :   left Expr, operator *Token, right Expr",
		"GroupingExpr :   expression Expr",
		"LiteralExpr  :   value interface{}",
		"UnaryExpr    :   operator *Token, right Expr",
	}

	defs := "package main\n\n"
	impl := ""

	for _, typeInfo := range types {

		parts := strings.Split(typeInfo, ":")
		name := strings.Trim(parts[0], " ")
		params := strings.Split(parts[1], ",")

		defs += fmt.Sprintf("type %s struct {\n", name)
		for _, param := range params {
			defs += fmt.Sprintf("  %s\n", strings.Trim(param, " "))
		}
		defs += "}\n\n"

		impl += fmt.Sprintf("func (expr *%s) accept(v ExprVisitor) Any {\n  return v.visit%s(expr)\n}\n\n", name, name)
	}

	os.Stdout.WriteString(defs)
	os.Stdout.WriteString(impl)
}
