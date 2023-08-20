package dsl

import "rexplaindsl/dev.yasint.rexplaindsl/api"
import "rexplaindsl/dev.yasint.rexplaindsl/complex"

func Either(expressions ...api.Expression) api.Expression {
	expression := ""
	for i, v := range expressions {
		expression += v.ToRegex()
		if i != len(expressions)-1 {
			expression += api.ALTERNATION
		}
	}
	return NonCaptureGroup(api.Expression{
		Construct: api.Operator,
		Instance:  nil,
		ToRegex: func() string {
			return expression
		},
	})
}

func EitherStr(strings ...string) api.Expression {
	trie := complex.NewTrieExpression()
	for _, v := range strings {
		trie.Insert(v)
	}
	return api.Expression{
		Construct: api.Operator,
		Instance:  nil,
		ToRegex: func() string {
			return trie.ToRegex()
		},
	}
}

func EitherStringsSet(strings []string) api.Expression {
	trie := &complex.TrieExpression{}
	trie.InsertAll(strings)
	return api.Expression{
		Construct: api.Operator,
		Instance:  nil,
		ToRegex: func() string {
			return trie.ToRegex()
		},
	}
}

func Concat(a api.Expression, b api.Expression) api.Expression {
	expression := a.ToRegex() + b.ToRegex()
	return api.Expression{
		Construct: api.Operator,
		Instance:  nil,
		ToRegex: func() string {
			return expression
		},
	}
}

func ConcatMultiple(expressions ...api.Expression) api.Expression {
	expression := ""
	for _, v := range expressions {
		expression += v.ToRegex()
	}
	return api.Expression{
		Construct: api.Operator,
		Instance:  nil,
		ToRegex: func() string {
			return expression
		},
	}
}
