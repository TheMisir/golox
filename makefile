GO       = go
AST_FILE = expression.go

generate:
	$(GO) run tool/generate_ast.go > $(AST_FILE)
	$(GO) fmt $(AST_FILE)

build:
	$(GO) build .

run:
	$(GO) run .