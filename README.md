# Monkey

[![Go Report Card](https://goreportcard.com/badge/github.com/tommymarto/Monkey)](https://goreportcard.com/report/github.com/tommymarto/Monkey) [![Build Status](https://dev.azure.com/tommymarto/Monkey/_apis/build/status/Monkey-Go%20(preview)-CI?repoName=tommymarto%2FMonkey&branchName=master)](https://dev.azure.com/tommymarto/Monkey/_build/latest?definitionId=2&repoName=tommymarto%2FMonkey&branchName=master)

![The official Monkey logo](https://interpreterbook.com/img/monkey_logo-d5171d15.png)

A [Monkey](https://monkeylang.org/) language AST-walking interpreter and a bytecode-compiled version (running in a custom VM) written in Go.
Built following the books [Writing an Interpreter in Go](https://interpreterbook.com/) and [Writing a Compiler in Go](https://compilerbook.com/).

The project features:

- a lexer with its own token definition and tokenization policies
- an AST producing parser implemented using the top-down Pratt approach
- a tree-walking evaluator with support for a syntactic macro system (Elixir-like).
- a compiler producing custom defined bytecode
- a custom stack-based VM capable of executing Monkey bytecode

All Monkey functionalities are available both in tree-walking and compiled mode.
Compiled Monkey is 3 to 4 times faster than interpreted Monkey!!!
You can run a sample benchmark (simply runs fib(35) and takes the execution time) typying `go run .\benchmark\ -engine=x` where x is either `vm` or `eval`

## What Monkey looks like

Variables (integers/booleans/strings), Arrays and Hash maps

```go
let x = (1 + 3) / 2 * 7  
let name = "Monkey"
let array = [x, name, true]
let dict = {name: 1, 2: x, true: array}
```

Functions

```go
let fibonacci = fn(x) {
    if (x == 0) {
        0                // Monkey supports implicit returning of values
    } else {
        if (x == 1) {
            return 1;      // ... and explicit return statements
        } else {
            fibonacci(x - 1) + fibonacci(x - 2); // and recursion!
        }
    }
}
```

Higher-order functions

```go
let map = fn(arr, f) {
    let iter = fn(arr, accumulated) {   // iter is a recursive closure!!
        if (len(arr) == 0) {
            accumulated
        } else {
            iter(rest(arr), push(accumulated, f(first(arr))));
        }
    };

    iter(arr, []);
};
```

function closures

```go
let add = fn(x) {
    return fn(y) { x + y }
};

let addTwo = add(2);

addTwo(4); // Output : 6
```

and macros

```go
let unless = macro(condition, consequence, alternative) {
    quote(if (!(unquote(condition))) {
            unquote(consequence);
        } else {
            unquote(alternative);
        });
};

unless(10 > 5, puts("not greater"), puts("greater")); // Output : "greater"
```

## The purpose of this project

This implementation of Monkey is for solely didactic purposes. The implementation makes heavy use of already existing Go objects for language-internal representation without any particular attention paid for optimization. The lexer ignores most of basic things like line numbers, the parser and the evaluator could be much more extended and the syntactic macro system severly lacks in error handling. With that said Monkey is easily extendable and Go garbage collector handles Monkey's garbage too!
A compiler and a VM has been added to the project. The shift from AST-walking to bytecode execution improved the performance by a factor from 3 to 4

### Try Monkey in the REPL

Simply type `go run .\main.go` *

You can choose the evaluation engine by specifying the flag `-engine=vm` for the compiled version or `-engine=eval` for the AST-walking version

You can view executing timing for each line directly in the REPL and make some nice comparison benchmark between engines.

\* REPL is currently implemented with no support for multiline statements/expressions. Might be added in future
