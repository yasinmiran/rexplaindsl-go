package dsl

import (
	"rexplaindsl/dev.yasint.rexplaindsl/api"
	"strconv"
)

type greedyQuantifier struct{}
type reluctantQuantifier struct{}

func OneOrMoreTimes(expression api.Expression) api.Expression {
	if isAGreedyQuantifier(expression) || isAReluctantQuantifier(expression) {
		panic("cannot apply + because it's already quantified")
	}
	return api.Expression{
		Construct: api.Repetition,
		Instance:  greedyQuantifier{},
		ToRegex: func() string {
			return NonCaptureGroup(expression).ToRegex() + api.PLUS
		},
	}
}

func ZeroOrMoreTimes(expression api.Expression) api.Expression {
	if isAGreedyQuantifier(expression) || isAReluctantQuantifier(expression) {
		panic("cannot apply * because it's already quantified")
	}
	return api.Expression{
		Construct: api.Repetition,
		Instance:  greedyQuantifier{},
		ToRegex: func() string {
			return NonCaptureGroup(expression).ToRegex() + api.ASTERISK
		},
	}
}

func ExactlyOrMoreTimes(times int, expression api.Expression) api.Expression {
	if isAGreedyQuantifier(expression) || isAReluctantQuantifier(expression) {
		panic("cannot apply {n,} because it's already quantified")
	}
	if times > 1000 {
		panic("max repetition is 1000")
	}
	if times == 0 {
		return ZeroOrMoreTimes(expression)
	}
	if times == 1 {
		return OneOrMoreTimes(expression)
	}
	return api.Expression{
		Construct: api.Repetition,
		Instance:  greedyQuantifier{},
		ToRegex: func() string {
			return NonCaptureGroup(expression).ToRegex() +
				api.OPEN_CURLY_BRACE + strconv.Itoa(times) +
				api.COMMA + api.CLOSE_CURLY_BRACE
		},
	}
}

func Optional(expression api.Expression) api.Expression {
	if isAGreedyQuantifier(expression) || isAReluctantQuantifier(expression) {
		panic("cannot apply ? because it's already quantified")
	}
	return api.Expression{
		Construct: api.Repetition,
		Instance:  greedyQuantifier{},
		ToRegex: func() string {
			return NonCaptureGroup(expression).ToRegex() +
				api.QUESTION_MARK
		},
	}
}

func Exactly(times int, expression api.Expression) api.Expression {
	if isAGreedyQuantifier(expression) || isAReluctantQuantifier(expression) {
		panic("cannot apply {n} because it's already quantified")
	}
	if times == 0 {
		panic("redundant sub-sequence")
	}
	if times == 1 {
		panic("redundant quantifier")
	}
	if times > 1000 {
		panic("max repetition is 1000")
	}
	return api.Expression{
		Construct: api.Repetition,
		Instance:  greedyQuantifier{},
		ToRegex: func() string {
			return NonCaptureGroup(expression).ToRegex() +
				api.OPEN_CURLY_BRACE + strconv.Itoa(times) +
				api.CLOSE_CURLY_BRACE
		},
	}
}

func Between(m int, n int, expression api.Expression) api.Expression {
	if isAGreedyQuantifier(expression) || isAReluctantQuantifier(expression) {
		panic("cannot apply {m,n} because it's already quantified")
	}
	if m > 1000 || n > 1000 {
		panic("max repetition is {1,1000}")
	}
	if m > n {
		panic("range is out of order")
	}
	if m == 0 && n == 0 {
		panic("redundant sub-sequence")
	} // below is default
	// Optimizations for the quantifiers
	if m == 0 && n == 1 {
		return Optional(expression)
	}
	if m == 1 && n == 1 {
		return expression
	}
	if m == n {
		return Exactly(m, expression)
	}
	return api.Expression{
		Construct: api.Repetition,
		Instance:  greedyQuantifier{},
		ToRegex: func() string {
			return NonCaptureGroup(expression).ToRegex() +
				api.OPEN_CURLY_BRACE + strconv.Itoa(m) + api.COMMA +
				strconv.Itoa(n) + api.CLOSE_CURLY_BRACE
		},
	}
}

func Lazy(expression api.Expression) api.Expression {
	if isAReluctantQuantifier(expression) {
		panic("already marked as lazy")
	}
	if !isAGreedyQuantifier(expression) {
		panic("must be a greedy 'quantifier'")
	}
	return api.Expression{
		Construct: api.Repetition,
		Instance:  greedyQuantifier{},
		ToRegex: func() string {
			return expression.ToRegex() + api.QUESTION_MARK
		},
	}
}

// go specific type checkers

func isAGreedyQuantifier(exp api.Expression) bool {
	if _, yes := exp.Instance.(greedyQuantifier); yes {
		return true
	}
	return false
}

func isAReluctantQuantifier(exp api.Expression) bool {
	if _, yes := exp.Instance.(reluctantQuantifier); yes {
		return true
	}
	return false
}
