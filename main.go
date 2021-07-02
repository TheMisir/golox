package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func runFromFile(name string) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}

	source := string(data)
	run(source)
}

func runFromStdin() {
	ctx := MakeContext()
	interpreter := MakeInterpreter(ctx)

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		ctx.hadError = false
		line := stdin.Text()

		scanner := MakeScanner(ctx, line)
		tokens := scanner.scanTokens()
		if ctx.hadError {
			continue
		}

		parser := MakeParser(ctx, tokens)
		statements, _ := parser.parse()
		if ctx.hadError {
			continue
		}

		resolver := MakeResolver(ctx, interpreter)
		resolver.resolve(statements)
		if ctx.hadError {
			continue
		}

		if len(statements) == 1 {
			switch val := statements[0].(type) {
			case *ExpressionStmt:
				result := interpreter.evaluate(val.expression)
				fmt.Fprintf(os.Stdout, "%v\n", result)
				break
			default:
				interpreter.interpret(statements)
				break
			}
		}
	}
}

func run(source string) {
	ctx := MakeContext()
	interpreter := MakeInterpreter(ctx)

	scanner := MakeScanner(ctx, source)
	scanner.scanTokens()
	if ctx.hadError {
		os.Exit(65)
	}

	parser := MakeParser(ctx, scanner.tokens)
	statements, _ := parser.parse()
	if !ctx.hadError {
		println("AST:")
		println(treePrinter.print(statements))

		os.Exit(65)
	}

	resolver := MakeResolver(ctx, interpreter)
	resolver.resolve(statements)
	if ctx.hadError {
		os.Exit(65)
	}

	interpreter.interpret(statements)
}

func main() {
	switch len(os.Args) {
	case 1:
		runFromStdin()
		break
	case 2:
		runFromFile(os.Args[1])
		break
	default:
		os.Stderr.WriteString("Syntax: golox [source]")
		os.Exit(64)
		break
	}
}
