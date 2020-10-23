// package repl
// basic Read Evaluate Print Loop. Handles one line inputs and prints its output in the terminal

package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
	"time"
)

const PROMPT = ">> "

func StartEval(in io.Reader, out io.Writer) {
	io.WriteString(out, "Running engine=eval\n")

	scanner := bufio.NewScanner(in)

	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		start := time.Now()
		result := evaluator.Eval(expanded, env)
		duration := time.Since(start)

		if result == nil {
			continue
		}

		io.WriteString(out, result.Inspect())
		io.WriteString(out, "\t\t")
		io.WriteString(out, duration.String())
		io.WriteString(out, "\n")
	}
}

func StartVM(in io.Reader, out io.Writer) {
	io.WriteString(out, "Running engine=vm\n")
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Whoops! Compilation failed:\n%s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)

		start := time.Now()
		err = machine.Run()
		duration := time.Since(start)
		if err != nil {
			fmt.Fprintf(out, "Whoops! Executing bytecode failed:\n%s\n", err)
			continue
		}

		result := machine.LastPoppedStackElem()
		if result == nil {
			continue
		}

		io.WriteString(out, result.Inspect())
		io.WriteString(out, "\t\t")
		io.WriteString(out, duration.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
