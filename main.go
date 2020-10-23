package main

import (
	"flag"
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

func main() {
	flag.Parse()

	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	switch *engine {
	case "vm":
		repl.StartVM(os.Stdin, os.Stdout)
	case "eval":
		repl.StartEval(os.Stdin, os.Stdout)
	default:
		fmt.Printf("Please specify a valid evaluation engine")
	}
}
