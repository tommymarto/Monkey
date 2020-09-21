# Monkey

![The official Monkey logo](https://interpreterbook.com/img/monkey_logo-d5171d15.png)

A [Monkey](https://monkeylang.org/) interpreter implementation in Go.
Built following the book [Writing an Interpreter in Go](https://interpreterbook.com/).

The project features:
- a lexer with its own token definition and tokenization policies
- an AST producing parser implemented using the top-down Pratt approach 
- a tree-walking evaluator with support for a syntactic macro system (Elixir-like).

## What Monkey looks like

Variables (integers/booleans/strings), Arrays and Hash maps
```
let x = (1 + 3) / 2 * 7  
let name = "Monkey"
let array = [x, name, true]
let dict = {name: 1, 2: x, true: array}
```

Functions
```
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
```
let map = fn(arr, f) {
	let iter = fn(arr, accumulated) {
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
```
let add = fn(x) {
	return fn(y) { x + y }
};

let addTwo = add(2);

addTwo(4); // Output : 6
```

and macros
```
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
This implementation of Monkey is for solely didactic purposes. The implementation makes heavy use of already existing Go objects for language-internal representation without any partitolar attention paid for optimization. The lexer ignores most of basic things like line numbers, the parser and the evaluator could be much more extended and the syntactic macro system severly lacks in error handling. With that said Monkey is easily extensible and Go garbage collector handles Monkey's garbage too!  

### Try Monkey in the REPL
Simply type `go run .\main.go` *

\* REPL is currently implemented with no support for multiline statements/expressions. Might be added in future
