GO       = go
AST_FILE = ast.go

ifeq ($(OS),Windows_NT)
	OS  = Win32
	EXT = .exe
else
	OS  = $(strip $(shell uname))
endif

OUT = build/$(OS)/golox$(EXT)

clean:
	$(RM) -rf build

generate:
	$(GO) run tool/generate_ast.go > $(AST_FILE)
	$(GO) fmt $(AST_FILE)

build:
	$(GO) build -o $(OUT)

run:
	$(GO) run .

all: clean generate build
