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
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		run(scanner.Text())
	}
}

func run(source string) {
	ctx := MakeContext()

	scanner := MakeScanner(ctx, source)
	scanner.scanTokens()

	parser := MakeParser(ctx, scanner.tokens)
	expression, err := parser.parse()

	if err != nil {
		return
	}

	println(MakeAstPrinter().print(expression))
	fmt.Printf("%v\n", MakeInterpereter().evaulate(expression))
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
