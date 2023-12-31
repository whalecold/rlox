package main

import "fmt"

type Environment struct {
	envs map[string]any
	// ancestor env
	enclosing *Environment
}

func NewEnvironmentWithAncestor(enclosing *Environment) *Environment {
	return &Environment{
		envs:      make(map[string]any),
		enclosing: enclosing,
	}
}

func NewEnvironment() *Environment {
	return &Environment{
		envs: make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.envs[name] = value
}

func (e *Environment) AssignAt(distance int, name *Token, value any) {
	e.ancestor(distance).envs[name.lexeme] = value
}

func (e *Environment) Assign(name *Token, value any) {
	if _, ok := e.envs[name.lexeme]; ok {
		e.envs[name.lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}
	Panic(name.line, fmt.Sprintf("Undefined variable '%s'.", name))
}

func (e *Environment) GetAt(distance int, name string) any {
	return e.ancestor(distance).envs[name]
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}

func (e *Environment) Get(name *Token) any {
	if val, ok := e.envs[name.lexeme]; ok {
		return val
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	Panic(name.line, fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
	return nil
}
