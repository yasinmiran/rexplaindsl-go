package complex

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
	"regexp"
	"rexplaindsl/dev.yasint.rexplaindsl/api"
	"rexplaindsl/dev.yasint.rexplaindsl/unicode"
	"rexplaindsl/dev.yasint.rexplaindsl/util"
)

type SetExpression struct {
	unicodeClasses *treeset.Set
	codePoints     *treeset.Set
	negated        bool
}

func NewSetExpression() SetExpression {
	return SetExpression{
		negated:        false,
		codePoints:     treeset.NewWith(utils.IntComparator),
		unicodeClasses: treeset.NewWith(utils.StringComparator),
	}
}

func (se SetExpression) Negate() {
	se.negated = true
}

func (se SetExpression) AddRange(codePointA int, codePointB int) {
	if util.IsValidCodePoint(codePointA) && util.IsValidCodePoint(codePointB) {
		if codePointA > codePointB {
			panic("character range is out of order")
		}
		if codePointA == codePointB {
			se.codePoints.Add(codePointA)
			return
		}
		for i := codePointA; i <= codePointB; i++ {
			se.codePoints.Add(i)
		}
	} else {
		panic("invalid codePoints")
	}
}

func (se SetExpression) AddChar(codePoint int) {
	if !util.IsValidCodePoint(codePoint) {
		panic("invalid codePoint")
	}
	se.codePoints.Add(codePoint)
}

func (se SetExpression) Union(b SetExpression) {
	if b.negated {
		se.codePoints.Remove(b.codePoints.Values()...)
	} else {
		se.codePoints.Add(b.codePoints.Values()...)
	}
}

func (se SetExpression) Intersection(b SetExpression) {
	if b.negated {
		se.codePoints.Remove(b.codePoints.Values()...)
	} else {
		modified := false
		iter := se.codePoints.Iterator()
		iter.Begin()
		for iter.Next() {
			val := iter.Value()
			if !b.codePoints.Contains(val) {
				se.codePoints.Remove(val)
				modified = true
			}
		}
		iter.End()
		fmt.Println("did modify reference set?", modified)
	}
}

func (se SetExpression) Difference(b SetExpression) {
	if !b.negated {
		b.Negate()
	}
	se.Intersection(b)
}

func (se SetExpression) WithUnicodeClass(block unicode.UnicodeScript, negated bool) {
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
	se.unicodeClasses.Add(expression)
}

func (se SetExpression) ToRegex() string {

	chars := toIntSlice(se.codePoints.Values())

	// return nothing if the set is empty
	if len(chars) == 0 && se.unicodeClasses.Empty() {
		return ""
	}
	// return only the unicode script class if it's a singleton
	if len(chars) == 0 {
		if se.unicodeClasses.Size() == 1 && !se.negated {
			return se.unicodeClasses.Values()[0].(string)
		}
	}

	if len(chars) == 1 && !se.negated && se.unicodeClasses.Empty() {
		return toRegexInterpretable(chars[0])
	}

	expression := "" + api.OPEN_SQUARE_BRACKET
	if se.negated {
		expression += api.CARAT
	}

	rangeStartIndex := -1
	isInRange := false

	for curIndex := range chars {
		if curIndex+1 < len(chars) {
			if (chars[curIndex+1] - chars[curIndex]) == 1 {
				if !isInRange {
					rangeStartIndex = curIndex
					isInRange = true
				}
				continue
			}
		}
		if isInRange {
			if (curIndex - rangeStartIndex) == 1 {
				expression += toRegexInterpretable(chars[rangeStartIndex]) +
					toRegexInterpretable(chars[curIndex])
			} else {
				expression += toRegexInterpretable(chars[rangeStartIndex]) +
					api.HYPHEN +
					toRegexInterpretable(chars[curIndex])
			}
			rangeStartIndex = -1
			isInRange = false
		} else {
			expression += toRegexInterpretable(chars[curIndex])
		}
	}

	if se.unicodeClasses.Size() > 0 {
		se.unicodeClasses.Each(func(index int, value interface{}) {
			expression += value.(string)
		})
	}

	expression += api.CLOSE_SQUARE_BRACKET
	return expression

}

func IsASetExpression(exp api.Expression) bool {
	if _, isSet := exp.Instance.(SetExpression); isSet {
		return true
	}
	return false
}

// private

func toIntSlice(any []interface{}) []int {
	intSlice := make([]int, len(any))
	for i := range any {
		intSlice[i] = any[i].(int)
	}
	return intSlice
}

func toRegexInterpretable(codePoint int) string {
	if util.IsISOControl(codePoint) || codePoint == 0x60 { // special case for tilde
		return fmt.Sprintf("\\x%02x", codePoint)
	}
	if util.IsSupplementaryCodePoint(codePoint) {
		return fmt.Sprintf("\\x{%x}", codePoint)
	}
	if util.IsBMPCodePoint(codePoint) {
		char := fmt.Sprintf("%c", rune(codePoint))
		var reserved = regexp.MustCompile(`[\\^\]/\-"'` + "`]")
		if match := reserved.MatchString(char); match {
			return api.BACKSLASH + char
		}
	}
	return fmt.Sprintf("%c", rune(codePoint)) // default case
}
