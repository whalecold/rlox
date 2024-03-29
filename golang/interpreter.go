package main

import "fmt"

type Interpreter struct {
	line    int
	env     *Environment
	globals *Environment
	locals  map[Expr]int
}

func NewInterpreter() *Interpreter {
	i := &Interpreter{
		globals: NewEnvironment(),
		locals:  make(map[Expr]int),
	}
	i.env = i.globals
	return injectPrimitives(i)
}

func (i *Interpreter) isTruthy(val any) bool {
	switch e := val.(type) {
	case bool:
		return e
	case nil:
		return false
	}
	return true
}

func (i *Interpreter) VisitUnaryExpr(expr Expr) any {
	e, ok := expr.(*Unary)
	if !ok {
		panic("should be unary type expr")
	}
	i.line = e.operator.line
	right := i.evaluate(e.right)
	switch e.operator.typ {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		return -i.float64Val(right)
	default:
		Panic(e.operator.line, fmt.Sprintf("Expect unary operator but got %v", e.operator.lexeme))
		return nil
	}
}

func (i *Interpreter) float64Val(in any) float64 {
	v, ok := in.(float64)
	if !ok {
		Panic(i.line, "should be float type")
	}
	return v
}

func (i *Interpreter) VisitBinaryExpr(expr Expr) any {
	e, ok := expr.(*Binary)
	if !ok {
		panic("should be binary type expr")
	}
	i.line = e.operator.line
	left := i.evaluate(e.left)
	right := i.evaluate(e.right)
	switch e.operator.typ {
	case MINUS:
		return i.float64Val(left) - i.float64Val(right)
	case SLASH:
		return i.float64Val(left) / i.float64Val(right)
	case STAR:
		return i.float64Val(left) * i.float64Val(right)
	case PLUS:
		switch l := left.(type) {
		case string:
			return l + fmt.Sprintf("%v", right)
		case float64:
			return l + i.float64Val(right)
		default:
			Panic(e.operator.line, fmt.Sprintf("Expect string or float type but got %v", l))
		}
	case GREATER:
		return i.float64Val(left) > i.float64Val(right)
	case GREATER_EQUAL:
		return i.float64Val(left) >= i.float64Val(right)
	case LESS:
		return i.float64Val(left) < i.float64Val(right)
	case LESS_EQUAL:
		return i.float64Val(left) <= i.float64Val(right)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	default:
		Panic(e.operator.line, fmt.Sprintf("Expect binary operator but got %v", e.operator.lexeme))
	}
	return nil
}

func (i *Interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) VisitLiteralExpr(expr Expr) any {
	e, ok := expr.(*Literal)
	if !ok {
		panic("should be literal type expr")
	}
	return e.value
}

func (i *Interpreter) VisitGroupingExpr(expr Expr) any {
	e, ok := expr.(*Grouping)
	if !ok {
		panic("should be grouping type expr")
	}
	return i.evaluate(e.expression)
}

func (i *Interpreter) VisitVariableExpr(expr Expr) any {
	e, ok := expr.(*Variable)
	if !ok {
		panic("should be variable type expr")
	}
	return i.lookupVariable(e.name, expr)
}

func (i *Interpreter) lookupVariable(name *Token, expr Expr) any {
	if distance, ok := i.locals[expr]; ok {
		return i.env.GetAt(distance, name.lexeme)
	}
	return i.globals.Get(name)
}

func (i *Interpreter) VisitCallExpr(expr Expr) any {
	e, ok := expr.(*Call)
	if !ok {
		panic("should be call type expr")
	}

	callee := i.evaluate(e.callee)
	args := make([]any, len(e.arguments))
	for k, v := range e.arguments {
		args[k] = i.evaluate(v)
	}
	function, ok := callee.(Callable)
	if !ok {
		Panic(e.paren.line, fmt.Sprintf("Expect callable but got %v", callee))
	}
	if len(args) != function.Arity() {
		Panic(e.paren.line, fmt.Sprintf("Expected %v arguments but got %v", function.Arity(), len(args)))
	}
	return function.Call(i, args)
}

func (i *Interpreter) VisitLogicalExpr(expr Expr) any {
	e, ok := expr.(*Logical)
	if !ok {
		panic("should be logical type expr")
	}
	left := i.evaluate(e.left)
	if i.isTruthy(left) {
		if e.operator.typ == OR {
			return left
		}
	} else {
		if e.operator.typ == AND {
			return left
		}
	}
	return i.evaluate(e.right)
}

func (i *Interpreter) VisitGetExpr(expr Expr) any {
	e, ok := expr.(*Get)
	if !ok {
		panic("should be get type expr")
	}
	val := i.evaluate(e.object)
	if o, ok := val.(*LoxInstance); ok {
		return o.Get(e.name)
	}
	Panic(e.name.line, "Only instances have properties")
	return nil
}

func (i *Interpreter) VisitSetExpr(expr Expr) any {
	e, ok := expr.(*Set)
	if !ok {
		panic("should be set type expr")
	}
	obj := i.evaluate(e.object)
	if o, ok := obj.(*LoxInstance); ok {
		val := i.evaluate(e.value)
		o.Set(e.name, val)
		return val
	}
	Panic(e.name.line, "Only instances have fields")
	return nil
}

func (i *Interpreter) VisitSuperExpr(expr Expr) any {
	e, ok := expr.(*Super)
	if !ok {
		panic("should be super type expr")
	}
	distance := i.locals[e]
	supperclass := i.env.GetAt(distance, "super").(*LoxClass)
	object := i.env.GetAt(distance-1, "this").(*LoxInstance)
	method := supperclass.FindMethod(e.method)
	if method == nil {
		Panic(e.method.line, fmt.Sprintf("Undefined property '%s'.", e.method.lexeme))
	}
	m, ok := method.(Callable)
	if !ok {
		Panic(e.method.line, fmt.Sprintf("'%s' is not a function.", e.method.lexeme))
	}
	return m.Bind(object)
}

func (i *Interpreter) VisitThisExpr(expr Expr) any {
	e, ok := expr.(*This)
	if !ok {
		panic("should be this type expr")
	}
	return i.lookupVariable(e.keyword, expr)
}

func (i *Interpreter) VisitAssignExpr(expr Expr) any {
	e, ok := expr.(*Assign)
	if !ok {
		panic("should be assign type")
	}
	// can't assign to undeclared variable
	val := i.evaluate(e.value)
	dis, ok := i.locals[expr]
	if ok {
		i.env.AssignAt(dis, e.name, val)
	} else {
		i.globals.Assign(e.name, val)
	}
	return val
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitClassStmt(stmt Stmt) any {
	s, ok := stmt.(*Class)
	if !ok {
		panic("should be class type stmt")
	}

	var superclass *LoxClass
	if s.superclass != nil {
		sc := i.evaluate(s.superclass)
		var ok bool
		superclass, ok = sc.(*LoxClass)
		if !ok {
			Panic(s.superclass.name.line, "Superclass must be a class")
		}
	}

	i.env.Define(s.name.lexeme, nil)

	if s.superclass != nil {
		i.env = NewEnvironmentWithAncestor(i.env)
		i.env.Define("super", superclass)
	}

	methods := make(map[string]Callable)
	for _, method := range s.methods {
		methods[method.name.lexeme] = NewCallable(method, i.env, method.name.lexeme == "init")
	}

	loxClass := NewLoxClass(s.name.lexeme, superclass, methods)

	if s.superclass != nil {
		i.env = i.env.enclosing
	}

	i.env.Assign(s.name, loxClass)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt Stmt) any {
	s, ok := stmt.(*Print)
	if !ok {
		panic("should be print type stmt")
	}
	val := i.evaluate(s.expr)
	fmt.Println(ToString(val))
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt Stmt) any {
	s, ok := stmt.(*Expression)
	if !ok {
		panic("should be expression type stmt")
	}
	return i.evaluate(s.expr)
}

func (i *Interpreter) VisitReturnStmt(stmt Stmt) any {
	s, ok := stmt.(*Return)
	if !ok {
		panic("should be return type stmt")
	}
	var val any
	if s.value != nil {
		val = i.evaluate(s.value)
	}
	panic(&ReturnPanic{Value: val})
}

func (i *Interpreter) VisitVarStmt(stmt Stmt) any {
	s, ok := stmt.(*Var)
	if !ok {
		panic("should be variable type stmt")
	}
	var val any
	if s.initializer != nil {
		val = i.evaluate(s.initializer)
	}
	i.env.Define(s.name.lexeme, val)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt Stmt) any {
	s, ok := stmt.(*If)
	if !ok {
		panic("should be if type stmt")
	}
	if i.isTruthy(i.evaluate(s.condition)) {
		i.execute(s.thenBranch)
	} else if s.elseBranch != nil {
		i.execute(s.elseBranch)
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt Stmt) any {
	s, ok := stmt.(*While)
	if !ok {
		panic("should be while type stmt")
	}
	for i.isTruthy(i.evaluate(s.condition)) {
		i.execute(s.body)
	}
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt Stmt) any {
	s, ok := stmt.(*Function)
	if !ok {
		panic("should be function type stmt")
	}
	i.env.Define(s.name.lexeme, NewCallable(s, i.env, false))
	return s
}

func (i *Interpreter) VisitBlockStmt(stmt Stmt) any {
	s, ok := stmt.(*Block)
	if !ok {
		panic("should be block type stmt")
	}
	i.executeBlock(s.statements, NewEnvironmentWithAncestor(i.env))
	return nil
}

func (i *Interpreter) execute(stmt Stmt) any {
	return stmt.Accept(i)
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) executeBlock(statements []Stmt, env *Environment) {
	previous := i.env
	i.env = env
	defer func() {
		i.env = previous
	}()
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) Execute(stmts []Stmt) any {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			hadError = true
		}
	}()
	for _, stmt := range stmts {
		i.execute(stmt)
	}
	return nil
}
