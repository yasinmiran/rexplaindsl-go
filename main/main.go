package main

import (
	"fmt"
	"rexplaindsl/dev.yasint.rexplaindsl/api"
	. "rexplaindsl/dev.yasint.rexplaindsl/dsl"
	"strings"
	"time"
)

func main() {
	dateMatchingExample()
	DoubleNumberMatchingWith1to3FractionDigits()
}

func dateMatchingExample() {

	months := []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul" /*"Aug",*/, "Sep", "Oct", "Nov", /*"Dec"*/
	}

	expression := api.NewReXPlainDSL(
		ExactLineMatch(
			IntegerRange(2012, time.Now().Year()),
			Literal("-"),
			EitherStr(months...),
			Literal("-"),
			Either(
				LeadingZero(IntegerRange(1, 9)),
				IntegerRange(10, 31),
			),
		),
	).Compile()

	fmt.Println(assertEquals("^(?:202[0-3]|201[2-9])\\-(?:Ap"+
		"r|Feb|J(?:an|u[ln])|Ma[ry]|Nov|Oct|Sep)\\-(?:(?:0?[1-9])"+
		"|(?:3[01]|[12][0-9]))$", expression.RegexExpression()))

}

func DoubleNumberMatchingWith1to3FractionDigits() {
	pattern := api.NewReXPlainDSL(
		ExactWordBoundary(
			IntegerRange(0, 1000),
			Literal("."),
			Between(1, 3, Digit()),
		),
	).Compile().PatternInstance()
	fmt.Println(pattern.String())
}

func assertEquals(expected string, actual string) bool {
	switch strings.Compare(expected, actual) {
	case -1:
		return false
	case 0:
		return true
	case 1:
		return false
	}
	return false
}
