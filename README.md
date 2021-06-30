# AST

Here's Lox expression grammar:

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

| Grammar notation | Code representation               |
| :--------------- | :-------------------------------- |
| Terminal         | Code to match and consume a token |
| Nonterminal      | Call to that rule’s function      |
| `\|`             | `if` or `switch` statement        |
| `*` or `+`       | `while` or `for` loop             |
| `?`              | `if` statement                    |
