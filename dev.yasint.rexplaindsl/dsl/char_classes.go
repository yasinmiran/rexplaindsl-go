package dsl

import (
	"rexplaindsl/dev.yasint.rexplaindsl/api"
	"rexplaindsl/dev.yasint.rexplaindsl/complex"
	"rexplaindsl/dev.yasint.rexplaindsl/unicode"
	"rexplaindsl/dev.yasint.rexplaindsl/util"
	"strings"
)

// Set construction functions.

func RangedSetStr(from string, to string) api.Expression {
	set := complex.NewSetExpression()
	set.AddRange(util.ToCodePoint(from), util.ToCodePoint(to))
	return api.Expression{
		Construct: api.Set,
		Instance:  set,
		ToRegex: func() string {
			return set.ToRegex()
		},
	}
}

func RangedSetCp(from int, to int) api.Expression {
	set := complex.NewSetExpression()
	set.AddRange(from, to)
	return api.Expression{
		Construct: api.Set,
		Instance:  set,
		ToRegex: func() string {
			return set.ToRegex()
		},
	}
}

func SimpleSetStr(chars ...string) api.Expression {
	set := complex.NewSetExpression()
	for _, v := range chars {
		set.AddChar(util.ToCodePoint(v))
	}
	return api.Expression{
		Construct: api.Set,
		Instance:  set,
		ToRegex: func() string {
			return set.ToRegex()
		},
	}
}

func SimpleSetCp(codePoints ...int) api.Expression {
	set := complex.NewSetExpression()
	for _, v := range codePoints {
		set.AddChar(v)
	}
	return api.Expression{
		Construct: api.Set,
		Instance:  set,
		ToRegex: func() string {
			return set.ToRegex()
		},
	}
}

func EmptySet() api.Expression {
	set := complex.NewSetExpression()
	return api.Expression{
		Construct: api.Set,
		Instance:  set,
		ToRegex: func() string {
			return set.ToRegex()
		},
	}
}

// Match any character construct.

func Anything() api.Expression {
	return SimpleSetStr(api.PERIOD)
}

// Set operator functions

func Negated(exp api.Expression) api.Expression {
	if complex.IsASetExpression(exp) {
		exp.Instance.(complex.SetExpression).Negate()
		return exp
	}
	panic("to negate it must be a set expression")
}

func Union(a api.Expression, b api.Expression) api.Expression {
	if complex.IsASetExpression(a) && complex.IsASetExpression(b) {
		a.Instance.(complex.SetExpression).
			Union(b.Instance.(complex.SetExpression))
		return a
	}
	panic("union only supported for set expressions")
}

func Difference(a api.Expression, b api.Expression) api.Expression {
	if complex.IsASetExpression(a) && complex.IsASetExpression(b) {
		a.Instance.(complex.SetExpression).
			Difference(b.Instance.(complex.SetExpression))
		return a
	}
	panic("union only supported for set expressions")
}

func Intersection(a api.Expression, b api.Expression) api.Expression {
	if complex.IsASetExpression(a) && complex.IsASetExpression(b) {
		a.Instance.(complex.SetExpression).
			Intersection(b.Instance.(complex.SetExpression))
		return a
	}
	panic("union only supported for set expressions")
}

func IncludeUnicodeScript(set api.Expression, block unicode.UnicodeScript, negated bool) api.Expression {
	if !complex.IsASetExpression(set) {
		panic("includeUnicodeScript only supported for set expressions")
	}
	set.Instance.(complex.SetExpression).
		WithUnicodeClass(block, negated)
	return set
}

// Posix

func LowerCase() api.Expression {
	return RangedSetStr("a", "z")
}

func UpperCase() api.Expression {
	return RangedSetStr("A", "Z")
}

func Ascii() api.Expression {
	return RangedSetCp(0x00, 0x7F)
}

func AsciiExtended() api.Expression {
	return RangedSetCp(0x00, 0xFF)
}

func Alphabetic() api.Expression {
	return Union(LowerCase(), UpperCase())
}

func Digit() api.Expression {
	return RangedSetStr("0", "9")
}

func NonDigit() api.Expression {
	return Negated(Digit())
}

func Alphanumeric() api.Expression {
	return Union(Alphabetic(), Digit())
}

func Punctuation() api.Expression {
	elements := "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	return SimpleSetStr(strings.Split(elements, "")...)
}

func Graphical() api.Expression {
	return Union(Alphanumeric(), Punctuation())
}

func Printable() api.Expression {
	return Union(Graphical(), SimpleSetCp(0x20))
}

func Blank() api.Expression {
	return SimpleSetCp(0x09, 0x20)
}

func HexDigit() api.Expression {
	return Union(
		RangedSetStr("A", "F"),
		Union(Digit(), RangedSetStr("a", "f")),
	)
}

func WhiteSpace() api.Expression {
	return SimpleSetCp(0x20, 0x9, 0xA, 0xB, 0xC, 0xD)
}

func NotWhiteSpace() api.Expression {
	return Negated(WhiteSpace())
}

func Word() api.Expression {
	return Union(Alphanumeric(), SimpleSetStr("_"))
}

func NotWord() api.Expression {
	return Negated(Word())
}

func Control() api.Expression {
	return Union(
		RangedSetCp(0x0, 0x1F),
		SimpleSetCp(0x7F),
	)
}

// Escape sequences

func Space() api.Expression {
	return SimpleSetStr(" ")
}

func BackSlash() api.Expression {
	return SimpleSetStr(`\\`)
}

func DoubleQuotes() api.Expression {
	return SimpleSetStr(`"`)
}

func SingleQuotes() api.Expression {
	return SimpleSetStr("'")
}

func BackTick() api.Expression {
	return SimpleSetCp(0x60)
}

func Bell() api.Expression {
	return SimpleSetCp(0x07)
}

func HorizontalTab() api.Expression {
	return SimpleSetCp(0x09)
}

func VerticalTab() api.Expression {
	return SimpleSetCp(0x0B)
}

func LineBreak() api.Expression {
	return SimpleSetCp(0x0A)
}

func FormFeed() api.Expression {
	return SimpleSetCp(0x0C)
}

func CarriageReturn() api.Expression {
	return SimpleSetCp(0x0D)
}
