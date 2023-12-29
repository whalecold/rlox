// Code generated by glox-gen. DO NOT EDIT.

package main

type Expr interface {
	Accept(Visitor) any
}

type Visitor interface {
	VisitBinaryExpr(Expr) any
	VisitGroupingExpr(Expr) any
	VisitLiteralExpr(Expr) any
	VisitUnaryExpr(Expr) any
}

type Binary struct {
	left     Expr
	operator *Token
	right    Expr
}

func (e *Binary) Accept(v Visitor) (ret any) {
	return v.VisitBinaryExpr(e)
}

type Grouping struct {
	expression Expr
}

func (e *Grouping) Accept(v Visitor) (ret any) {
	return v.VisitGroupingExpr(e)
}

type Literal struct {
	value any
}

func (e *Literal) Accept(v Visitor) (ret any) {
	return v.VisitLiteralExpr(e)
}

type Unary struct {
	operator *Token
	right    Expr
}

func (e *Unary) Accept(v Visitor) (ret any) {
	return v.VisitUnaryExpr(e)
}