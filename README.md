# GoLox

The [lox language](https://craftinginterpreters.com/the-lox-language.html) interpreter written in Golang.

---

Check [improved](https://github.com/TheMisir/golox/tree/improved) branch 
for an improved version of the interpreter.

## Build

If you have [Go](https://golang.org) installed, you can build the 
project with:

```sh
go build .
```

Or you can use make to build the project:

```sh
make build
```

## Notes

These are notes I took during reading the book that helped me to 
write the code.

### Grammar

Here's Lox AST node grammar. Those are simple notations that 
contains everything needed for implementing parser for given 
statements.

```plain
program        → declaration* EOF ;

declaration    → classDecl
               | funDecl
               | varDecl
               | statement ;

classDecl      → "class" IDENTIFIER ( "<" IDENTIFIER )?
                 "{" function* "}" ;

funDecl        → "fun" function ;
function       → IDENTIFIER "(" parameters? ")" block ;
parameters     → IDENTIFIER ( "," IDENTIFIER )* ;

statement      → exprStmt
               | forStmt
               | ifStmt
               | printStmt
               | returnStmt
               | whileStmt
               | block ;

returnStmt     → "return" expression? ";" ;

forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
                 expression? ";"
                 expression? ")" statement ;

whileStmt      → "while" "(" expression ")" statement ;

ifStmt         → "if" "(" expression ")" statement
               ( "else" statement )? ;

block          → "{" declaration* "}" ;

exprStmt       → expression ";" ;

printStmt      → "print" expression ";" ;

varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;

expression     → assignment ;
assignment     → ( call "." )? IDENTIFIER "=" assignment
               | logic_or ;
logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary | call ;
call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
arguments      → expression ( "," expression )* ;
primary        → "true" | "false" | "nil" | "this"
               | NUMBER | STRING | IDENTIFIER | "(" expression ")"
               | "super" "." IDENTIFIER ;
```

And the body of the rule translates to code roughly like:

| Grammar notation    | Code representation               |
| :------------------ | :-------------------------------- |
| Terminal            | Code to match and consume a token |
| Nonterminal         | Call to that rule’s function      |
| <code>&#124;</code> | `if` or `switch` statement        |
| `*` or `+`          | `while` or `for` loop             |
| `?`                 | `if` statement                    |
