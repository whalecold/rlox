package main

import (
	"bufio"
	"fmt"
	"os"
)

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
		if hadError {
			return
		}
	}
}

func run(content string) {
	scanner := NewScanner(content)
	parser := &Parser{tokens: scanner.ScanTokens()}
	stmts := parser.ParseStmts()
	if stmts == nil {
		return
	}
	//fmt.Println(expr.Accept(&AstPrinter{}))
	inter := &Interpreter{}

	inter.Interpreter(stmts)

	//for _, token := range scanner.ScanTokens() {
	//	fmt.Println("token: ", token)
	//}
}

func toString(in any) string {
	switch in.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, in)
	case int:
		return fmt.Sprintf("%d", in)
	default:
		return fmt.Sprintf("%v", in)
	}
}

func interpret(expr Expr) any {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			hadError = true
		}
	}()
	return toString(expr.Accept(&Interpreter{}))
}
