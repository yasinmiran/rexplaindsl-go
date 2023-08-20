package dsl

import (
	"rexplaindsl/dev.yasint.rexplaindsl/api"
	"rexplaindsl/dev.yasint.rexplaindsl/complex"
	"strconv"
)

func LeadingZero(expression api.Expression) api.Expression {
	return NonCaptureGroup(
		api.Expression{
			Construct: api.Numeric,
			Instance:  nil,
			ToRegex: func() string {
				return "0" + api.QUESTION_MARK + expression.ToRegex()
			},
		},
	)
}

func IntegerRange(from int, to int) api.Expression {
	if from > to {
		panic("integer range is out of order")
	}
	if from == to {
		return Literal(strconv.Itoa(from))
	}
	if from >= 0 && to <= 9 {
		return RangedSetStr(strconv.Itoa(from), strconv.Itoa(to))
	}
	return NonCaptureGroup(
		api.Expression{
			Construct: api.Numeric,
			Instance:  nil,
			ToRegex: func() string {
				return complex.NewRangeExpression(from, to).ToRegex()
			},
		},
	)
}
