package main

import "fmt"

type LoxInstance struct {
	loxClass *LoxClass
	fileds   map[string]any
}

func NewLoxInstance(loxClass *LoxClass) *LoxInstance {
	return &LoxInstance{
		loxClass: loxClass,
		fileds:   make(map[string]any),
	}
}

func (lox *LoxInstance) ToString() string {
	return lox.loxClass.name + " instance"
}

func (lox *LoxInstance) Set(name *Token, value any) {
	lox.fileds[name.lexeme] = value
}

func (lox *LoxInstance) Get(name *Token) any {
	if val, ok := lox.fileds[name.lexeme]; ok {
		return val
	}

	method := lox.loxClass.FindMethod(name)
	if method != nil {
		if fn, ok := method.(Callable); ok {
			return fn.Bind(lox)
		} else {
			Panic(name.line, fmt.Sprintf("'%s' is not a function.", name.lexeme))
		}
	}

	Panic(name.line, fmt.Sprintf("Undefined property '%s'.", name))
	return nil
}
