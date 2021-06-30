package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

var printAst bool = true

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
	printer := MakeAstPrinter()

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

		if printAst {
			fmt.Fprintf(os.Stdout, "AST:\n")
			for _, statement := range statements {
				fmt.Fprintf(os.Stdout, "%s\n", printer.printStmt(statement))
			}
			fmt.Fprintf(os.Stdout, "\n\n")
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

	printer := MakeAstPrinter()
	scanner := MakeScanner(ctx, source)
	scanner.scanTokens()
	if ctx.hadError {
		os.Exit(65)
	}

	parser := MakeParser(ctx, scanner.tokens)
	statements, _ := parser.parse()
	if ctx.hadError {
		os.Exit(65)
	}

	if printAst {
		fmt.Fprintf(os.Stdout, "AST:\n")
		for _, statement := range statements {
			fmt.Fprintf(os.Stdout, "%s\n", printer.printStmt(statement))
		}
		fmt.Fprintf(os.Stdout, "\n\n")
	}

	interpreter := MakeInterpreter(ctx)
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
