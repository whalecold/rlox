package main

type LoxClass struct {
	name       string
	methods    map[string]Callable
	superclass *LoxClass
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]Callable) *LoxClass {
	return &LoxClass{
		name:       name,
		methods:    methods,
		superclass: superclass,
	}
}

func (lc *LoxClass) Arity() int {
	return 0
}
func (lc *LoxClass) Call(*Interpreter, []any) any {
	return NewLoxInstance(lc)
}

func (lc *LoxClass) ToString() string {
	return lc.name
}

func (lc *LoxClass) Bind(*LoxInstance) Callable {
	return lc
}

func (lc *LoxClass) FindMethod(name *Token) any {
	if val, ok := lc.methods[name.lexeme]; ok {
		return val
	}
	if lc.superclass != nil {
		return lc.superclass.FindMethod(name)
	}
	return nil
}
