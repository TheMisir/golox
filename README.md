# GoLox

The [lox language](https://craftinginterpreters.com/the-lox-language.html) interpreter written in Golang.

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


## Improvements

This branch adds a few "unoriginal" features to the original Lox 
implementation.



#### Functions are now considered as expressions

Unlike original Lox which considered function definitions as statements,
this branch considers function definitions as a expression that resolves
to function itself (instance of LoxFunction in runtime). This allows
directly providing function as a value to another function.

```lox
fun call(callback) {
  return callback();
}

call(fun greet() {  // Here we can directly declare a function as a value
  return "Hello";
});
```

#### Anonymous functions

This is not a huge thing, but it is a nice feature that allows to write
functions without name which could be used as lambdas for calling another
functions. For example we can write previous code like that:

```lox
fun call(callback) {
  return callback();
}

call(fun () {  // We don't have to provide a name for the function
  return "Hello";
});
```

I also modified the resolver to detect some edge cases like declaring 
nameless functions as a class method.

#### Modified `for` statement handling

With modified lox for statements are not de-sugarized into a while
statement. This is done to allow for writing more complex loops and
supporting `break` and `continue` statements.

#### Loop breaking and continuing

Now you can use `break` and `continue` statements in loops - both while
and for loops.


Syntax is simple as below:

```lox
while (true) {
  if (condition) {
    break;
  }
}
```

#### Including files

Lox now supports including files. This is useful for writing libraries
that can be used in other projects.

```lox
include "lib/math.lox";
```

Include statement can be chained with other statements like `if` statement.
That's useful when you want to include files only when some condition is 
met.

```lox
if (isWindows) {
  include "lib/windows.lox";
} else {
  include "lib/unix.lox";
}
```

#### What's Next?

I'm planning to add more features to the language and to the interpreter.
Here is a list of things I want to add:

 - [ ] Support for `namespace` blocks
 - [ ] Array literals
 - [ ] Map literals, _maybe_
 - [ ] Add more operations to standard library
 - [ ] ~~Foreign function calls to dynamic libraries~~

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
               | breakStmt
               | continueStmt
               | includeStmt
               | block ;

includeStmt    → "include" STRING ";" ;

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
               | "super" "." IDENTIFIER
               | funDecl ;
```

And the body of the rule translates to code roughly like:

| Grammar notation    | Code representation               |
| :------------------ | :-------------------------------- |
| Terminal            | Code to match and consume a token |
| Nonterminal         | Call to that rule’s function      |
| <code>&#124;</code> | `if` or `switch` statement        |
| `*` or `+`          | `while` or `for` loop             |
| `?`                 | `if` statement                    |
