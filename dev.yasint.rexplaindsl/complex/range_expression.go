package complex

import (
	"rexplaindsl/dev.yasint.rexplaindsl/api"
	"strconv"
)

type RangeExpression struct {
	Start int
	End   int
}

func NewRangeExpression(start int, end int) RangeExpression {
	return RangeExpression{
		Start: start,
		End:   end,
	}
}

func leftBounds(start int, end int) []partialRange {
	var result []partialRange
	for start < end {
		ran := fromStart(start)
		result = append(result, ran)
		start = ran.end + 1
	}
	return result
}

func rightBounds(start int, end int) []partialRange {
	var result []partialRange
	for start < end {
		ran := fromEnd(end)
		result = append(result, ran)
		end = ran.start - 1
	}
	return reverseRange(result)
}

func reverseRange(elements []partialRange) []partialRange {
	var reversed []partialRange
	for i := range elements {
		n := elements[len(elements)-1-i]
		reversed = append(reversed, n)
	}
	return reversed
}

func (re RangeExpression) ToRegex() string {

	left := leftBounds(re.Start, re.End)
	lastLeft := left[len(left)-1]
	left = left[:len(left)-1]

	right := rightBounds(lastLeft.start, re.End)
	firstRight := right[0]
	right = right[1:]

	var merged []partialRange
	for _, v := range left {
		merged = append(merged, v)
	}
	if !lastLeft.overlaps(firstRight) {
		merged = append(merged, lastLeft)
		merged = append(merged, firstRight)
	} else {
		merged = append(merged, join(lastLeft, firstRight))
	}
	for _, v := range right {
		merged = append(merged, v)
	}

	expression := ""
	for i, e := range reverseRange(merged) {
		expression += e.ToRegex()
		if i+1 != len(merged) {
			expression += api.ALTERNATION
		}

	}

	return expression

}

// Partial range expression. To multiple ranges

type partialRange struct {
	start int
	end   int
}

func (r partialRange) overlaps(another partialRange) bool {
	return r.end > another.start && another.end > r.start
}

func fromEnd(end int) partialRange {
	chars := []rune(strconv.Itoa(end))
	for i := len(chars) - 1; i >= 0; i-- {
		if chars[i] == '9' {
			chars[i] = '0'
		} else {
			chars[i] = '0'
			break
		}
	}
	parsed, _ := strconv.Atoi(string(chars))
	return partialRange{parsed, end}
}

func fromStart(start int) partialRange {
	chars := []rune(strconv.Itoa(start))
	for i := len(chars) - 1; i >= 0; i-- {
		if chars[i] == '0' {
			chars[i] = '9'
		} else {
			chars[i] = '9'
			break
		}
	}
	parsed, _ := strconv.Atoi(string(chars))
	return partialRange{start, parsed}
}

func join(a partialRange, b partialRange) partialRange {
	return partialRange{a.start, b.end}
}

func (r partialRange) ToRegex() string {

	startString := strconv.Itoa(r.start)
	endString := strconv.Itoa(r.end)
	expression := ""
	repeatedCount := 0
	prevDigitA, prevDigitB := -1, -1

	for pos := 0; pos < len(startString); pos++ {
		curDigitA, _ := strconv.Atoi(string(startString[pos]))
		curDigitB, _ := strconv.Atoi(string(endString[pos]))
		if curDigitA == curDigitB {
			expression += strconv.Itoa(curDigitA)
		} else {
			if prevDigitA == curDigitA && prevDigitB == curDigitB {
				repeatedCount++
				if !(pos == len(startString)-1) {
					continue
				} else {
					repeatedCount++ // Signifies that we're in a repetition construct
					expression += api.OPEN_CURLY_BRACE + strconv.Itoa(repeatedCount) + api.CLOSE_CURLY_BRACE
					break
				}
			}
			if repeatedCount > 0 {
				expression += api.OPEN_CURLY_BRACE + strconv.Itoa(repeatedCount) + api.CLOSE_CURLY_BRACE
				repeatedCount = 0 // Reset the repetition construct to continue new elements
			}
			expression += api.OPEN_SQUARE_BRACKET + strconv.Itoa(curDigitA)
			if !(curDigitB-curDigitA == 1) {
				expression += api.HYPHEN
			}
			expression += strconv.Itoa(curDigitB) + api.CLOSE_SQUARE_BRACKET
			prevDigitA = curDigitA
			prevDigitB = curDigitB
		}
	}

	return expression

}
