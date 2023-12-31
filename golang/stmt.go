// Code generated by glox-gen. DO NOT EDIT.

package main

type Stmt interface {
	Accept(StmtVisitor) any
}

type StmtVisitor interface {
	VisitExpressionStmt(Stmt) any
	VisitFunctionStmt(Stmt) any
	VisitPrintStmt(Stmt) any
	VisitReturnStmt(Stmt) any
	VisitVarStmt(Stmt) any
	VisitBlockStmt(Stmt) any
	VisitIfStmt(Stmt) any
	VisitWhileStmt(Stmt) any
}

type Expression struct {
	expr Expr
}

func (e *Expression) Accept(v StmtVisitor) (ret any) {
	return v.VisitExpressionStmt(e)
}

type Function struct {
	name   *Token
	params []*Token
	body   []Stmt
}

func (e *Function) Accept(v StmtVisitor) (ret any) {
	return v.VisitFunctionStmt(e)
}

type Print struct {
	expr Expr
}

func (e *Print) Accept(v StmtVisitor) (ret any) {
	return v.VisitPrintStmt(e)
}

type Return struct {
	keyword *Token
	value   Expr
}

func (e *Return) Accept(v StmtVisitor) (ret any) {
	return v.VisitReturnStmt(e)
}

type Var struct {
	name        *Token
	initializer Expr
}

func (e *Var) Accept(v StmtVisitor) (ret any) {
	return v.VisitVarStmt(e)
}

type Block struct {
	statements []Stmt
}

func (e *Block) Accept(v StmtVisitor) (ret any) {
	return v.VisitBlockStmt(e)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (e *If) Accept(v StmtVisitor) (ret any) {
	return v.VisitIfStmt(e)
}

type While struct {
	condition Expr
	body      Stmt
}

func (e *While) Accept(v StmtVisitor) (ret any) {
	return v.VisitWhileStmt(e)
}
