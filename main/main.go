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

	fmt.Println(assertEquals("^(?:2020|201[2-9])\\-(?:Apr|Feb|J(?:an|u[ln])|Ma[ry]|"+
		"Nov|Oct|Sep)\\-(?:(?:0?[1-9])|(?:3[01]|[12][0-9]))$", expression.RegexExpression()))

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
