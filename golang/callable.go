package main

import (
	"fmt"
	"time"
)

type ReturnPanic struct {
	Value any
}

type Callable interface {
	Arity() int
	Call(*Interpreter, []any) any
	ToString() string
}

type CallableFunc func([]any) any

type callableImpl struct {
	argsNumber  int
	fn          CallableFunc
	primitive   bool
	declaration *Function
	closure     *Environment
}

func NewPrimitive(arity int, fn CallableFunc) Callable {
	return Callable(&callableImpl{
		argsNumber: arity,
		primitive:  true,
		fn:         fn,
	})
}

func NewCallable(declaration *Function, e *Environment) Callable {
	return Callable(&callableImpl{
		primitive:   false,
		declaration: declaration,
		closure:     e,
	})
}

func (c *callableImpl) ToString() string {
	if c.primitive {
		return "<fn primitive>"
	}
	return "<fn " + c.declaration.name.lexeme + ">"
}

func (c *callableImpl) Call(i *Interpreter, args []any) (ret any) {
	if len(args) != c.Arity() {
		panic(fmt.Sprintf("Expected %v arguments but got %v", c.argsNumber, len(args)))
	}
	if c.primitive {
		return c.fn(args)
	}
	e := NewEnvironmentWithAncestor(c.closure)
	for k, v := range c.declaration.params {
		e.Define(v.lexeme, args[k])
	}
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *ReturnPanic:
				ret = r.(*ReturnPanic).Value
			default:
				panic(r)
			}
		}
	}()
	i.executeBlock(c.declaration.body, e)
	return nil
}

func (c *callableImpl) Arity() int {
	if c.primitive {
		return c.argsNumber
	}
	return len(c.declaration.params)
}

func injectPrimitives(i *Interpreter) *Interpreter {
	i.env.Define("clock", NewPrimitive(0, func(args []any) any {
		return time.Now().Unix()
	}))
	return i
}
