package main

import (
	"fmt"
	"os"
)

type LoxContext struct {
	hasError bool
}

func MakeContext() *LoxContext {
	return &LoxContext{
		hasError: false,
	}
}

func (c *LoxContext) error(line int, message string, a ...interface{}) {
	c.report(line, "", message)
}

func (c *LoxContext) tokenError(token *Token, message string) {
	if token.tokenType == EOF {
		c.report(token.line, " at end", message)
	} else {
		c.report(token.line, fmt.Sprintf(" at '%s'", token.lexme), message)
	}
}

func (c *LoxContext) report(line int, where string, message string, a ...interface{}) {
	c.hasError = true
	message = fmt.Sprintf(message, a...)
	fmt.Fprintf(os.Stderr, "[line %v] Error%s: %s\n", line, where, message)
}
