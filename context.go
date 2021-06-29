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
	os.Stderr.WriteString(fmt.Sprintf(message+fmt.Sprintf(" (line: %v)\n", line), a...))
	c.hasError = true
}
