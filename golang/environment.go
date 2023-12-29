package main

import "fmt"

type Environment struct {
	envs map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		envs: make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.envs[name] = value
}

func (e *Environment) Get(name *Token) any {
	if _, ok := e.envs[name.lexeme]; !ok {
		Panic(name.line, fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
	}
	return e.envs[name.lexeme]
}
