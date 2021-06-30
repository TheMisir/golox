package main

import (
	"bufio"
	"io/ioutil"
	"os"
)

var ctx = MakeContext()
var interpreter = MakeInterpreter(ctx)

func runFromFile(name string) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}

	source := string(data)
	run(source, true)
}

func runFromStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ctx.hadError = false
		run(scanner.Text(), true)
	}
}

func run(source string, continueOnError bool) {
	scanner := MakeScanner(ctx, source)
	scanner.scanTokens()
	if ctx.hadError {
		if continueOnError {
			return
		} else {
			os.Exit(65)
		}
	}

	parser := MakeParser(ctx, scanner.tokens)
	statements, _ := parser.parse()
	if ctx.hadError {
		if continueOnError {
			return
		} else {
			os.Exit(65)
		}
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
