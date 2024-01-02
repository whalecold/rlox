package main

import "fmt"

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

type ClassType int

const (
	CLASSNONE ClassType = iota
	CLASSCLASS
	CLASSSUB
)

type Resolver struct {
	inter        *Interpreter
	scopes       []map[string]bool
	currnetFunc  FunctionType
	currentClass ClassType
}

func NewResolver(inter *Interpreter) *Resolver {
	return &Resolver{
		inter:        inter,
		scopes:       []map[string]bool{},
		currnetFunc:  NONE,
		currentClass: CLASSNONE,
	}
}

func (r *Resolver) VisitBlockStmt(stmt Stmt) any {
	s, ok := stmt.(*Block)
	if !ok {
		panic("should be block type stmt")
	}
	r.beginScope()
	r.resolveStmts(s.statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitVarStmt(stmt Stmt) any {
	s, ok := stmt.(*Var)
	if !ok {
		panic("should be variable type stmt")
	}
	r.declare(s.name)
	if s.initializer != nil {
		r.resolveExpr(s.initializer)
	}
	r.define(s.name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr Expr) any {
	e, ok := expr.(*Assign)
	if !ok {
		panic("should be assign type expr")
	}
	r.resolveExpr(e.value)
	r.resolveLocal(e, e.name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr Expr) any {
	e, ok := expr.(*Binary)
	if !ok {
		panic("should be binary type expr")
	}
	r.resolveExpr(e.left)
	r.resolveExpr(e.right)
	return nil
}

func (r *Resolver) VisitSetExpr(expr Expr) any {
	e, ok := expr.(*Set)
	if !ok {
		panic("should be set type expr")
	}
	r.resolveExpr(e.object)
	r.resolveExpr(e.value)
	return nil
}

func (r *Resolver) VisitSuperExpr(expr Expr) any {
	e, ok := expr.(*Super)
	if !ok {
		panic("should be super type expr")
	}
	if r.currentClass == CLASSNONE {
		Panic(e.keyword.line, "Can't use 'super' outside of a class.")
	}
	if r.currentClass != CLASSSUB {
		Panic(e.keyword.line, "Can't use 'super' in a class with no superclass.")
	}
	r.resolveLocal(e, e.keyword)
	return nil
}

func (r *Resolver) VisitGetExpr(expr Expr) any {
	e, ok := expr.(*Get)
	if !ok {
		panic("should be get type expr")
	}
	r.resolveExpr(e.object)
	return nil
}

func (r *Resolver) VisitCallExpr(expr Expr) any {
	e, ok := expr.(*Call)
	if !ok {
		panic("should be call type expr")
	}
	r.resolveExpr(e.callee)
	for _, arg := range e.arguments {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr Expr) any {
	e, ok := expr.(*Grouping)
	if !ok {
		panic("should be grouping type expr")
	}
	r.resolveExpr(e.expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr Expr) any {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr Expr) any {
	e, ok := expr.(*Logical)
	if !ok {
		panic("should be logical type expr")
	}
	r.resolveExpr(e.left)
	r.resolveExpr(e.right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr Expr) any {
	e, ok := expr.(*Unary)
	if !ok {
		panic("should be unary type expr")
	}
	r.resolveExpr(e.right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr Expr) any {
	e, ok := expr.(*Variable)
	if !ok {
		panic("should be variable type expr")
	}
	if len(r.scopes) != 0 {
		if val, ok := r.scopes[len(r.scopes)-1][e.name.lexeme]; ok && !val {
			Panic(e.name.line, "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(e, e.name)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt Stmt) any {
	s, ok := stmt.(*Function)
	if !ok {
		panic("should be function type stmt")
	}
	r.declare(s.name)
	r.define(s.name)
	r.resolveFunction(s, FUNCTION)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt Stmt) any {
	s, ok := stmt.(*Expression)
	if !ok {
		panic("should be expression type stmt")
	}
	r.resolveExpr(s.expr)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt Stmt) any {
	s, ok := stmt.(*If)
	if !ok {
		panic("should be if type stmt")
	}
	r.resolveExpr(s.condition)
	r.resolveStmt(s.thenBranch)
	if s.elseBranch != nil {
		r.resolveStmt(s.elseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt Stmt) any {
	s, ok := stmt.(*Print)
	if !ok {
		panic("should be print type stmt")
	}
	r.resolveExpr(s.expr)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt Stmt) any {
	s, ok := stmt.(*Return)
	if !ok {
		panic("should be return type stmt")
	}
	if r.currnetFunc == NONE {
		Panic(s.keyword.line, "Can't return from top-level code.")
	}
	if s.value != nil {
		if r.currnetFunc == INITIALIZER {
			Panic(s.keyword.line, "Can't return from initializer.")
		}
		r.resolveExpr(s.value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt Stmt) any {
	s, ok := stmt.(*While)
	if !ok {
		panic("should be while type stmt")
	}
	r.resolveExpr(s.condition)
	r.resolveStmt(s.body)
	return nil
}

func (r *Resolver) resolveFunction(fn *Function, typ FunctionType) {
	enclosingFunction := r.currnetFunc
	r.currnetFunc = typ
	defer func() {
		r.currnetFunc = enclosingFunction
	}()

	r.beginScope()
	for _, param := range fn.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(fn.body)
	r.endScope()
}

func (r *Resolver) define(token *Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[token.lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if r.scopes[i][name.lexeme] {
			r.inter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) declare(token *Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[token.lexeme]; ok {
		Panic(token.line, "Already a variable with this name in this scope.")
	}
	scope[token.lexeme] = false
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) VisitThisExpr(expr Expr) any {
	e, ok := expr.(*This)
	if !ok {
		panic("should be this type expr")
	}
	if r.currentClass == CLASSNONE {
		Panic(e.keyword.line, "Can't use 'this' outside of a class.")
	}
	r.resolveLocal(expr, e.keyword)
	return nil
}

func (r *Resolver) resolveExprs(exprs []Expr) {
	for _, expr := range exprs {
		expr.Accept(r)
	}
}

func (r *Resolver) Resolve(stmts []Stmt) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			hadError = true
		}
	}()
	r.resolveStmts(stmts)
}

func (r *Resolver) VisitClassStmt(stmt Stmt) any {
	s, ok := stmt.(*Class)
	if !ok {
		panic("should be class type stmt")
	}
	enclosingClass := r.currentClass
	r.currentClass = CLASSCLASS

	r.declare(s.name)

	if s.superclass != nil && s.name.lexeme == s.superclass.name.lexeme {
		Panic(s.superclass.name.line, "A class can't inherit from itself.")
	}

	if s.superclass != nil {
		r.currentClass = CLASSSUB
		r.resolveExpr(s.superclass)
	}

	if s.superclass != nil {
		r.beginScope()
		r.scopes[len(r.scopes)-1]["super"] = true
	}

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range s.methods {
		functionType := METHOD
		if method.name.lexeme == "init" {
			functionType = INITIALIZER
		}
		r.resolveFunction(method, functionType)
	}
	r.endScope()

	if s.superclass != nil {
		r.endScope()
	}

	r.define(s.name)

	r.currentClass = enclosingClass

	return nil
}

func (r *Resolver) resolveStmts(stmts []Stmt) {
	for _, stmt := range stmts {
		stmt.Accept(r)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}
func (r *Resolver) endScope() {
	if len(r.scopes) == 0 {
		panic("unreachable")
	}
	r.scopes = r.scopes[:len(r.scopes)-1]
}
