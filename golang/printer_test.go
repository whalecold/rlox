package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstPrinter(t *testing.T) {
	out := &Literal{
		value: 45.67,
	}
	assert.Equal(t, 45.67, out.Accept(&AstPrinter{}))
	expr := &Binary{
		left: &Unary{
			operator: Token{
				typ:    MINUS,
				lexeme: "-",
				line:   1,
			},
			right: &Literal{
				value: 123,
			},
		},
		operator: &Token{
			typ:    STAR,
			lexeme: "*",
			line:   1,
		},
		right: &Grouping{
			expression: &Literal{
				value: 45.67,
			},
		},
	}
	assert.Equal(t, "(* (- 123) (group 45.67))", expr.Accept(&AstPrinter{}))
}
