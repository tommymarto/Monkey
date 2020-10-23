package vm

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, tests)
}

func TestBooleanExpression(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}

	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}

	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Null},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}

	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}",
			map[object.HashKey]int64{},
		},
		{
			"{1: 2, 3: 4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 3}).HashKey(): 4,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			`
				let fivePlusTen = fn() { 5 + 10; };
				fivePlusTen();
			`,
			15,
		},
		{
			`
				let one = fn() { 1; };
				let two = fn() { 2; };
				one() + two();
			`,
			3,
		},
		{
			`
				let a = fn() { 1; };
				let b = fn() { a() + 1; };
				let c = fn() { b() + 1; };
				c();
			`,
			3,
		},
		{
			`
				let earlyExit = fn() { return 99; 100; };
				earlyExit();
			`,
			99,
		},
		{
			`
				let earlyExit = fn() { return 99; return 100; };
				earlyExit();
			`,
			99,
		},
		{
			`
				let noReturn = fn() { };
				noReturn();
			`,
			Null,
		},
		{
			`
				let noReturn = fn() { };
				let noReturnTwo = fn() { noReturn(); };
				noReturn();
				noReturnTwo();
			`,
			Null,
		},
	}

	runVmTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			`
				let returnsOne = fn() { 1; };
				let returnsOneReturner = fn() { returnsOne; };
				returnsOneReturner()()
			`,
			1,
		},
		{
			`
				let returnsOneReturner = fn() { 
					let returnsOne = fn() { 1; };
					returnsOne; 
				};
				returnsOneReturner()()
			`,
			1,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			`
				let one = fn() { let one = 1; one; };
				one();
			`,
			1,
		},
		{
			`
				let oneAndTwo = fn() { let one = 1; let two = 2; one + two; }
				oneAndTwo();
			`,
			3,
		},
		{
			`
				let oneAndTwo = fn() { let one = 1; let two = 2; one + two; }
				let threeAndFour = fn() { let three = 3; let four = 4; three + four; }
				oneAndTwo() + threeAndFour();
			`,
			10,
		},
		{
			`
				let firstFoobar = fn() { let foobar = 50; foobar; }
				let secondFoobar = fn() { let foobar = 100; foobar; }
				firstFoobar() + secondFoobar();
			`,
			150,
		},
		{
			`
				let globalSeed = 50;
				let minusOne = fn() { let num = 1; globalSeed - num; };
				let minusTwo = fn() { let num = 2; globalSeed - num; };
				minusOne() + minusTwo();
			`,
			97,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			`
				let identity = fn(a) { a };
				identity(4);
			`,
			4,
		},
		{
			`
				let sum = fn(a, b) { a + b; }
				sum(1, 2);
			`,
			3,
		},
		{
			`
				let sum = fn(a, b) { 
					let c = a + b;
					c; 
				};
				sum(1, 2);
			`,
			3,
		},
		{
			`
				let sum = fn(a, b) { 
					let c = a + b;
					c; 
				};
				sum(1, 2) + sum(3, 4);
			`,
			10,
		},
		{
			`
				let sum = fn(a, b) { 
					let c = a + b;
					c; 
				};
				let outer = fn() {
					sum(1, 2) + sum(3, 4);
				}
				outer();
			`,
			10,
		},
		{
			`
				let globalNum = 10;
				let sum = fn(a, b) { 
					let c = a + b;
					c + globalNum; 
				};
				let outer = fn() {
					sum(1, 2) + sum(3, 4) + globalNum;
				}
				outer() + globalNum;
			`,
			50,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			`fn() { 1; }(1);`,
			`wrong number of arguments: want=0, got=1`,
		},
		{
			`fn(a) { a; }();`,
			`wrong number of arguments: want=1, got=0`,
		},
		{
			`fn(a, b) { a + b; }(1);`,
			`wrong number of arguments: want=2, got=1`,
		},
	}

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error but resulted in none.")
		}

		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error: want=%q, got=%q", tt.expected, err)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world!")`, 12},
		{
			`
				len(1)
			`,
			&object.Error{
				Message: "argument to `len` not supported, got INTEGER",
			},
		},
		{
			`
				len("one", "two")
			`,
			&object.Error{
				Message: "wrong number of arguments. got=2, want=1",
			},
		},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, Null},
		{`first([1, 2, 3])`, 1},
		{`first([])`, Null},
		{
			`
				first(1)
			`,
			&object.Error{
				Message: "argument to `first` must be ARRAY, got INTEGER",
			},
		},
		{`last([1, 2, 3])`, 3},
		{`last([])`, Null},
		{
			`
				last(1)
			`,
			&object.Error{
				Message: "argument to `last` must be ARRAY, got INTEGER",
			},
		},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, Null},
		{`push([], 1)`, []int{1}},
		{
			`
				push(1, 1)
			`,
			&object.Error{
				Message: "argument to `push` must be ARRAY, got INTEGER",
			},
		},
	}

	runVmTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			`
				let newClosure = fn(a) {
					fn() { a }
				};
				let closure = newClosure(99);
				closure()
			`,
			99,
		},
		{
			`
				let newAdder = fn(a, b) {
					fn(c) { a + b + c }
				};
				let adder = newAdder(1, 2);
				adder(8)
			`,
			11,
		},
		{
			`
				let newAdder = fn(a, b) {
					let c = a + b;
					fn(d) { c + d }
				};
				let adder = newAdder(1, 2);
				adder(8)
			`,
			11,
		},
		{
			`
				let newAdderOuter = fn(a, b) {
					let c = a + b;
					fn(d) { 
						let e = c + d
						fn(f) { e + f }
					}
				};
				let newAdderInner = newAdderOuter(1, 2);
				let adder = newAdderInner(3)
				adder(8)
			`,
			14,
		},
		{
			`
				let a = 1
				let newAdderOuter = fn(b) {
					fn(c) {
						fn(d) { a + b + c + d }
					}
				};
				let newAdderInner = newAdderOuter(2);
				let adder = newAdderInner(3)
				adder(8)
			`,
			14,
		},
		{
			`
				let newClosure = fn(a, b) {
					let one = fn() { a; }
					let two = fn() { b; }
					fn() { one() + two(); }
				};
				let closure = newClosure(9, 90);
				closure()
			`,
			99,
		},
	}

	runVmTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			`
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						countDown(x - 1);
					}
				};
				countDown(1)
			`,
			0,
		},
		{
			`
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						countDown(x - 1);
					}
				};
				let wrapper = fn() {
					countDown(1);
				}
				wrapper()
			`,
			0,
		},
		{
			`
				let wrapper = fn() {
					let countDown = fn(x) {
						if (x == 0) {
							return 0;
						} else {
							countDown(x - 1);
						}
					};
					countDown(1)
				}
				wrapper()
			`,
			0,
		},
	}

	runVmTests(t, tests)
}

func TestRecursiveFibonacci(t *testing.T) {
	tests := []vmTestCase{
		{
			`
				let fib = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						if (x == 1) {
							return 1;
						} else {
							fib(x-1) + fib(x-2)
						}
					}
				}
				fib(15)
			`,
			610,
		},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		for i, constant := range comp.Bytecode().Constants {
			fmt.Printf("CONSTANT %d %p (%T):\n", i, constant, constant)

			switch constant := constant.(type) {
			case *object.CompiledFunction:
				fmt.Printf(" Instructions:\n%s", constant.Instructions)
			case *object.Integer:
				fmt.Printf(" Value: %d\n", constant.Value)
			}
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []int:
		err := testArrayObject(expected, actual)
		if err != nil {
			t.Errorf("testArrayObject failed: %s", err)
		}
	case map[object.HashKey]int64:
		err := testHashObject(expected, actual)
		if err != nil {
			t.Errorf("testHashObject failed: %s", err)
		}
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("object is not Error: %T (%+v)", actual, actual)
			return
		}
		if errObj.Message != expected.Message {
			t.Errorf("wrong error message. expected=%q, got=%q", expected.Message, errObj.Message)
		}
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}

	return nil
}

func testArrayObject(expected []int, actual object.Object) error {
	result, ok := actual.(*object.Array)
	if !ok {
		return fmt.Errorf("object is not Array. got=%T (%+v)", actual, actual)
	}

	if len(result.Elements) != len(expected) {
		return fmt.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(result.Elements))
	}

	for i, expEl := range expected {
		err := testIntegerObject(int64(expEl), result.Elements[i])
		if err != nil {
			return fmt.Errorf("test testIntegerObject failed: %s", err)
		}
	}

	return nil
}

func testHashObject(expected map[object.HashKey]int64, actual object.Object) error {
	result, ok := actual.(*object.Hash)
	if !ok {
		return fmt.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
	}

	if len(result.Pairs) != len(expected) {
		return fmt.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(result.Pairs))
	}

	for expK, expV := range expected {
		pair, ok := result.Pairs[expK]
		if !ok {
			return fmt.Errorf("no pair for given key in Pairs")
		}

		err := testIntegerObject(expV, pair.Value)
		if err != nil {
			return fmt.Errorf("testIntegerObject failed: %s", err)
		}
	}

	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
