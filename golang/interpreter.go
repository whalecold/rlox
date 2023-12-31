package main

import "fmt"

type Interpreter struct {
	line    int
	env     *Environment
	globals *Environment
}

func NewInterpreter() *Interpreter {
	i := &Interpreter{
		globals: NewEnvironment(),
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
	return i.env.Get(e.name)
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

func (i *Interpreter) VisitAssignExpr(expr Expr) any {
	e, ok := expr.(*Assign)
	if !ok {
		panic("should be assign type")
	}
	// can't assign to undeclared variable
	val := i.evaluate(e.value)
	i.env.Assign(e.name, val)
	return val
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
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
	i.env.Define(s.name.lexeme, NewCallable(s, i.env))
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

func (i *Interpreter) Interpreter(stmts []Stmt) any {
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
