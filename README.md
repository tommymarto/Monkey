# Monkey

![The official Monkey logo](https://interpreterbook.com/img/monkey_logo-d5171d15.png)

A [Monkey](https://monkeylang.org/) interpreter implementation in Go.
Built following the book [Writing an Interpreter in Go](https://interpreterbook.com/).

The project features a lexer, an AST producing parser and a tree-walking evaluator.

## What Monkey looks like

Basics
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

and closures
```
let add = fn(x) {
  return fn(y) { x + y }
};

let addTwo = add(2);

addTwo(4); // Output : 6
```

### Try Monkey in the REPL
Simply type `go run .\main.go`
