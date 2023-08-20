package api

import (
	"regexp"
)

type Flag string

const (
	CaseInsensitive Flag = "i"
	Multiline       Flag = "m"
	DotAll          Flag = "s"
)

type Framework interface {
	GetMatchedGroups() interface{}
	Compile(flags ...Flag) ReXPlainDSL
	PatternInstance() *regexp.Regexp
	RegexExpression() string
}

type ReXPlainDSL struct {
	expression string
	pattern    *regexp.Regexp
}

func NewReXPlainDSL(expressions ...Expression) ReXPlainDSL {
	expr := ""
	for _, v := range expressions {
		expr += v.ToRegex()
	}
	return ReXPlainDSL{
		expression: expr,
		pattern:    nil,
	}
}

func (rs ReXPlainDSL) GetMatchedGroups() interface{} {
	return nil
}

func (rs ReXPlainDSL) Compile(flags ...Flag) ReXPlainDSL {
	if len(flags) == 0 {
		rs.pattern = regexp.MustCompile(rs.expression)
		return rs
	}
	prefix := "("
	for _, f := range flags {
		prefix += string(f)
	}
	prefix += ")"
	rs.pattern = regexp.MustCompile(prefix + rs.expression)
	return rs
}

func (rs ReXPlainDSL) PatternInstance() *regexp.Regexp {
	if rs.pattern == nil {
		panic("pattern instance is null! call Compile(Flag...)")
	}
	return rs.pattern
}

func (rs ReXPlainDSL) RegexExpression() string {
	return rs.expression
}
