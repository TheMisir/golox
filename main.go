package main

import (
	"bufio"
	"fmt"
	"os"
)

func runFromStdin() {
	ctx := MakeContext()
	interpreter := MakeInterpreter(ctx)

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		ctx.hadError = false
		line := stdin.Text()

		source := &Source{Name: "<stdin>", Code: line}
		scanner := MakeScanner(ctx, source)
		tokens := scanner.scanTokens()
		if ctx.hadError {
			continue
		}

		parser := MakeParser(ctx, tokens)
		statements, _ := parser.parse()
		if ctx.hadError {
			continue
		}

		sourceResolver := MakeFileSourceResolver("")
		resolver := MakeResolver(ctx, interpreter, sourceResolver)
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

func runFromFile(name string) {
	ctx := MakeContext()
	interpreter := MakeInterpreter(ctx)
	sourceResolver := MakeFileSourceResolver("")

	source, err := sourceResolver.Resolve(ctx, name)
	if err != nil {
		panic(err)
	}

	if ctx.hadError {
		println("AST:")
		println(treePrinter.print(source.Body))

		os.Exit(65)
	}

	resolver := MakeResolver(ctx, interpreter, sourceResolver)
	resolver.resolve(source.Body)
	if ctx.hadError {
		os.Exit(65)
	}

	interpreter.interpret(source.Body)
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
