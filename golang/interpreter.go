package main

import "fmt"

type Interpreter struct {
	line int
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
		panic("should be unary type")
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
		panic("should be binary type")
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
		panic("should be literal type")
	}
	return e.value
}

func (i *Interpreter) VisitGroupingExpr(expr Expr) any {
	e, ok := expr.(*Grouping)
	if !ok {
		panic("should be grouping type")
	}
	return i.evaluate(e.expression)
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitErrorExprExpr(expr Expr) any {
	return expr.Accept(i)
}
