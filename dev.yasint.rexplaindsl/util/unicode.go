package util

// Go specific unicode helpers

const minSupplementaryCodePoint = 0x010000
const maxCodePoint = 0x10FFFF

func IsBMPCodePoint(codepoint int) bool {
	return codepoint>>16 == 0
}

func IsSupplementaryCodePoint(codePoint int) bool {
	return codePoint >= minSupplementaryCodePoint && codePoint < maxCodePoint+1
}

func IsValidCodePoint(codePoint int) bool {
	plane := codePoint >> 16
	return plane < ((maxCodePoint + 1) >> 16)
}

func IsISOControl(codePoint int) bool {
	return codePoint <= 0x9F && (codePoint >= 0x7F || (codePoint>>5 == 0))
}
