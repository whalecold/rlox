package main

import (
	"fmt"
	"os"
	"strings"
)

type Ast struct {
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		return
	}
	outputDir := os.Args[1]
	defineAst(outputDir, "Expr", []string{
		"Binary:left Expr,operator *Token,right Expr",
		"Grouping:expression Expr",
		"Literal:value any",
		"Unary:operator *Token,right Expr",
	},
	)
}

func defineAst(outputDir, baseName string, types []string) {
	path := outputDir + "/" + baseName + ".go"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	file.Write([]byte("// Code generated by glox-gen. DO NOT EDIT.\n\n"))
	file.Write([]byte("package main\n\n"))

	// expr interface
	file.Write([]byte("type " + baseName + " interface {\n"))
	file.Write([]byte("  Accept(Visitor) any \n"))
	file.Write([]byte("}\n\n"))

	// visitor interface
	file.Write([]byte("type Visitor interface {\n"))
	for _, t := range types {
		parts := strings.Split(t, ":")
		class := strings.TrimSpace(parts[0])
		file.Write([]byte("  Visit" + class + "Expr(" + baseName + ") any \n"))
	}
	file.Write([]byte("}\n\n"))

	for _, t := range types {
		parts := strings.Split(t, ":")
		class := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])
		defineType(file, baseName, class, fields)
	}
	file.Close()
}

var recoverStr = `
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			ret = nil
		}
	}()

`

func defineType(file *os.File, baseName, className, fields string) {
	file.Write([]byte("type " + className + " struct {\n"))
	fieldsList := strings.Split(fields, ",")
	for _, field := range fieldsList {
		fieldParts := strings.Split(field, " ")
		fieldName := strings.TrimSpace(fieldParts[0])
		fieldType := strings.TrimSpace(fieldParts[1])
		file.Write([]byte("\t" + fieldName + " " + fieldType + "\n"))
	}
	file.Write([]byte("}\n\n"))

	// Accept function
	file.Write([]byte("func (e *" + className + ") Accept(v Visitor) (ret any) {\n"))
	file.Write([]byte("  return v.Visit" + className + "Expr(e)\n"))
	file.Write([]byte("}\n\n"))
}
