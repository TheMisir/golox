package main

import (
	"bufio"
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

	for _, token := range scanner.tokens {
		println(token.toString())
	}
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
