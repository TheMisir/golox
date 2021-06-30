# GoLox

The [lox language](https://craftinginterpreters.com/the-lox-language.html) interpreter written in Golang.

## Notes

These are notes I took during reading the book that helped me to 
write the code.

### Expressions

Here's Lox expression grammar. Those are simple notations that 
contains everything needed for implementing parser for given 
statements.

```plain
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
```

And the body of the rule translates to code roughly like:

| Grammar notation    | Code representation               |
| :------------------ | :-------------------------------- |
| Terminal            | Code to match and consume a token |
| Nonterminal         | Call to that rule’s function      |
| <code>&#124;</code> | `if` or `switch` statement        |
| `*` or `+`          | `while` or `for` loop             |
| `?`                 | `if` statement                    |
