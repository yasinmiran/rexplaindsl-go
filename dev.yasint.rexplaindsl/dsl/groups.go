package dsl

import (
	"rexplaindsl/dev.yasint.rexplaindsl/api"
	"rexplaindsl/dev.yasint.rexplaindsl/util"
)

func NonCaptureGroup(expressions ...api.Expression) api.Expression {
	expression := api.PAREN_OPEN + api.QUESTION_MARK + api.COLON
	for _, v := range expressions {
		expression += v.ToRegex()
	}
	expression += api.PAREN_CLOSE
	return api.Expression{
		Construct: api.Group,
		Instance:  nil,
		ToRegex: func() string {
			return expression
		},
	}
}

func NamedCaptureGroup(name string, expressions ...api.Expression) api.Expression {
	expression := api.PAREN_OPEN + api.QUESTION_MARK + api.NAMED_CAPTURE_GROUP_PREFIX + api.LESS_THAN
	expression += util.AsRegexGroupName(name) + api.GREATER_THAN
	for _, v := range expressions {
		expression += v.ToRegex()
	}
	expression += api.PAREN_CLOSE
	return api.Expression{
		Construct: api.Group,
		Instance:  nil,
		ToRegex: func() string {
			return expression
		},
	}
}

func CaptureGroup(expressions ...api.Expression) api.Expression {
	expression := api.PAREN_OPEN
	for _, v := range expressions {
		expression += v.ToRegex()
	}
	expression += api.PAREN_CLOSE
	return api.Expression{
		Construct: api.Group,
		Instance:  nil,
		ToRegex: func() string {
			return expression
		},
	}
}
