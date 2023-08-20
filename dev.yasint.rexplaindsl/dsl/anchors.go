package dsl

import "rexplaindsl/dev.yasint.rexplaindsl/api"

func WordBoundary() api.Expression {
	return api.Expression{
		Construct: api.Anchor,
		Instance:  nil,
		ToRegex: func() string {
			return api.WORD_BOUNDARY
		},
	}
}

func NonWordBoundary() api.Expression {
	return api.Expression{
		Construct: api.Anchor,
		Instance:  nil,
		ToRegex: func() string {
			return api.NON_WORD_BOUNDARY
		},
	}
}

func StartOfLine() api.Expression {
	return api.Expression{
		Construct: api.Anchor,
		Instance:  nil,
		ToRegex: func() string {
			return api.CARAT
		},
	}
}

func EndOfLine(crlf bool) api.Expression {
	return api.Expression{
		Construct: api.Anchor,
		Instance:  nil,
		ToRegex: func() string {
			if crlf {
				return "\x0D?" + api.DOLLAR
			}
			return api.DOLLAR
		},
	}
}

func StartOfText() api.Expression {
	return api.Expression{
		Construct: api.Anchor,
		Instance:  nil,
		ToRegex: func() string {
			return api.BEGINNING_OF_TEXT
		},
	}
}

func EndOfText() api.Expression {
	return api.Expression{
		Construct: api.Anchor,
		Instance:  nil,
		ToRegex: func() string {
			return api.END_OF_TEXT
		},
	}
}

func ExactLineMatch(expressions ...api.Expression) api.Expression {
	return api.Expression{
		Construct: api.Anchor,
		Instance:  nil,
		ToRegex: func() string {
			var exp = api.CARAT
			for _, se := range expressions {
				exp += se.ToRegex()
			}
			exp += api.DOLLAR
			return exp
		},
	}
}

func ExactWordBoundary(expressions ...api.Expression) api.Expression {
	return api.Expression{
		Construct: api.Anchor,
		Instance:  nil,
		ToRegex: func() string {
			var exp = api.WORD_BOUNDARY
			for _, se := range expressions {
				exp += se.ToRegex()
			}
			exp += api.WORD_BOUNDARY
			return exp
		},
	}
}
