package main

import "fmt"

var hadError bool

func Panic(line int, message any) {
	hadError = true
	panic(fmt.Sprintf("[line %v] message: %v\n", line, message))
}

func Error(line int, message any) {
	report("Error ", line, "", message)
}

func Debug(line int, message any) {
	report("Debug ", line, "", message)
}

func report(prefix string, line int, where, message any) {
	fmt.Printf("%s[line %v] where: %v message: %v\n", prefix, line, where, message)
	hadError = true
}
