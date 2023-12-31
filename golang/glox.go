package main

import (
	"bufio"
	"fmt"
	"os"
)

var inter = NewInterpreter()

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		return
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	run(string(content))

}

func runPrompt() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">")
		scanner.Scan()
		text := scanner.Text()
		if text == "exit" || text == "q" {
			break
		}
		run(text)
	}
}

func run(content string) {
	scanner := NewScanner(content)
	parser := &Parser{tokens: scanner.ScanTokens()}
	stmts := parser.ParseStmts()
	if hadError {
		return
	}
	if stmts == nil {
		return
	}
	resolver := NewResolver(inter)
	resolver.Resolve(stmts)
	if hadError {
		return
	}
	inter.Execute(stmts)
}

func ToString(in any) string {
	if in == nil {
		return "nil"
	}
	switch in.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, in)
	case int:
		return fmt.Sprintf("%d", in)
	case Callable:
		return in.(Callable).ToString()
	default:
		return fmt.Sprintf("%v", in)
	}
}
