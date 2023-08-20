package util

import (
	"fmt"
	"regexp"
)

var Reserved = regexp.MustCompile(`[<(\[{\\^\-=$!|\]})?*+.>/]`)
var ValidGroupName = regexp.MustCompile(`^[^[:punct:][:digit:][:space:]]\w{1,15}$`)

func AsRegexLiteral(someString string) string {
	if len(someString) == 1 { // quickly return if its only 1 char
		if match := Reserved.MatchString(someString); match {
			return "\\" + someString
		}
	}
	codePoints := []rune(someString)
	composed := ""
	for _, v := range codePoints {
		if IsSupplementaryCodePoint(int(v)) {
			composed += "\\x{" + fmt.Sprintf("%x", int(v)) + "}"
			// in this implementation we don't need to skip the next rune
		} else {
			s := fmt.Sprintf("%c", v)
			if match := Reserved.MatchString(s); match {
				composed += "\\" + s
			} else {
				composed += s
			}
		}
	}
	return composed
}

func AsRegexGroupName(name string) string {
	valid := ValidGroupName.MatchString(name)
	if !valid {
		panic("invalid capture group name")
	}
	return name
}

func ToCodePoint(someString string) int {
	var points = []rune(someString)
	if len(points) != 1 {
		panic("expected a bmp or astral symbol")
	}
	return int(points[0])
}
