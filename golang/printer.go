package main

//type AstPrinter struct{}
//
//func parenthesize(name string, expr ...any) string {
//	var b strings.Builder
//	b.Grow(len(name) + 5*len(expr))
//	b.WriteString("(")
//	b.WriteString(name)
//	for _, e := range expr {
//		b.WriteString(" ")
//		b.WriteString(fmt.Sprintf("%v", e))
//	}
//	b.WriteString(")")
//	return b.String()
//}
//
//func (ap *AstPrinter) VisitUnaryExpr(e Expr) any {
//	unary, ok := e.(*Unary)
//	if !ok {
//		panic("should be unary type")
//	}
//	return parenthesize(unary.operator.lexeme, unary.right.Accept(ap))
//}
//
//func (ap *AstPrinter) VisitGetExpr(e Expr) any {
//	return nil
//}
//
//func (ap *AstPrinter) VisitSetExpr(e Expr) any {
//	return nil
//}
//
//func (ap *AstPrinter) VisitBinaryExpr(e Expr) any {
//	binary, ok := e.(*Binary)
//	if !ok {
//		panic("should be binary type")
//	}
//
//	return parenthesize(binary.operator.lexeme, binary.left.Accept(ap), binary.right.Accept(ap))
//}
//
//func (ap *AstPrinter) VisitCallExpr(e Expr) any {
//	return nil
//}
//
//func (ap *AstPrinter) VisitAssignExpr(e Expr) any {
//	return nil
//}
//
//func (ap *AstPrinter) VisitGroupingExpr(e Expr) any {
//	grouping, ok := e.(*Grouping)
//	if !ok {
//		panic("should be grouping type")
//	}
//	return parenthesize("group", grouping.expression.Accept(ap))
//}
//
//func (ap *AstPrinter) VisitLiteralExpr(e Expr) any {
//	literal, ok := e.(*Literal)
//	if !ok {
//		panic("should be literal type")
//	}
//	return literal.value
//}
//
//func (ap *AstPrinter) VisitVariableExpr(e Expr) any {
//	return nil
//}
//
//func (ap *AstPrinter) VisitLogicalExpr(e Expr) any {
//	return nil
//}
