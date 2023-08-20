package dsl

import (
	"rexplaindsl/dev.yasint.rexplaindsl/api"
	"rexplaindsl/dev.yasint.rexplaindsl/unicode"
	"rexplaindsl/dev.yasint.rexplaindsl/util"
)

func Literal(someString string) api.Expression {
	return api.Expression{
		Construct: api.Literal,
		Instance:  nil,
		ToRegex: func() string {
			return util.AsRegexLiteral(someString)
		},
	}
}

func QuotedLiteral(someString string) api.Expression {
	return api.Expression{
		Construct: api.Literal,
		Instance:  nil,
		ToRegex: func() string {
			exp := api.QUOTE_START + someString + api.QUOTE_END
			return exp
		},
	}
}

func UnicodeScriptLiteral(block unicode.UnicodeScript, negated bool) api.Expression {
	expression := ""
	if len(block) == 1 {
		if negated {
			expression += "\\P" + string(block)
		} else {
			expression += "\\p" + string(block)
		}
	} else {
		if negated {
			expression += "\\P" + api.OPEN_CURLY_BRACE + string(block) + api.CLOSE_CURLY_BRACE
		} else {
			expression += "\\p" + api.OPEN_CURLY_BRACE + string(block) + api.CLOSE_CURLY_BRACE
		}
	}
	return api.Expression{
		Construct: api.Literal,
		Instance:  nil,
		ToRegex: func() string {
			return expression
		},
	}
}
