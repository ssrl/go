// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unicode

// Bit masks for each code point under U+0100, for fast lookup.
const (
	pC  = 1 << iota // a control character.
	pP              // a punctuation character.
	pN              // a numeral.
	pS              // a symbolic character.
	pZ              // a spacing character.
	pLu             // an upper-case letter.
	pLl             // a lower-case letter.
	pp              // a printable character according to Go's definition.
	pg  = pp | pZ   // a graphical character according to the Unicode definition.
)

// GraphicRanges defines the set of graphic characters according to Unicode.
var GraphicRanges = []*RangeTable{
	L, M, N, P, S, Zs,
}

// PrintRanges defines the set of printable characters according to Go.
// ASCII space, U+0020, is handled separately.
var PrintRanges = []*RangeTable{
	L, M, N, P, S,
}

// IsGraphic reports whether the rune is defined as a Graphic by Unicode.
// Such characters include letters, marks, numbers, punctuation, symbols, and
// spaces, from categories L, M, N, P, S, Zs.
func IsGraphic(rune int) bool {
	// We cast to uint32 to avoid the extra test for negative,
	// and in the index we cast to uint8 to avoid the range check.
	if uint32(rune) <= MaxLatin1 {
		return properties[uint8(rune)]&pg != 0
	}
	return IsOneOf(GraphicRanges, rune)
}

// IsPrint reports whether the rune is defined as printable by Go. Such
// characters include letters, marks, numbers, punctuation, symbols, and the
// ASCII space character, from categories L, M, N, P, S and the ASCII space
// character.  This categorization is the same as IsGraphic except that the
// only spacing character is ASCII space, U+0020.
func IsPrint(rune int) bool {
	if uint32(rune) <= MaxLatin1 {
		return properties[uint8(rune)]&pp != 0
	}
	return IsOneOf(PrintRanges, rune)
}

// IsOneOf reports whether the rune is a member of one of the ranges.
// The rune is known to be above Latin-1.
func IsOneOf(set []*RangeTable, rune int) bool {
	for _, inside := range set {
		if Is(inside, rune) {
			return true
		}
	}
	return false
}

// IsControl reports whether the rune is a control character.
// The C (Other) Unicode category includes more code points
// such as surrogates; use Is(C, rune) to test for them.
func IsControl(rune int) bool {
	if uint32(rune) <= MaxLatin1 {
		return properties[uint8(rune)]&pC != 0
	}
	// All control characters are < Latin1Max.
	return false
}

// IsLetter reports whether the rune is a letter (category L).
func IsLetter(rune int) bool {
	if uint32(rune) <= MaxLatin1 {
		return properties[uint8(rune)]&(pLu|pLl) != 0
	}
	return Is(Letter, rune)
}

// IsMark reports whether the rune is a mark character (category M).
func IsMark(rune int) bool {
	// There are no mark characters in Latin-1.
	return Is(Mark, rune)
}

// IsNumber reports whether the rune is a number (category N).
func IsNumber(rune int) bool {
	if uint32(rune) <= MaxLatin1 {
		return properties[uint8(rune)]&pN != 0
	}
	return Is(Number, rune)
}

// IsPunct reports whether the rune is a Unicode punctuation character
// (category P).
func IsPunct(rune int) bool {
	if uint32(rune) <= MaxLatin1 {
		return properties[uint8(rune)]&pP != 0
	}
	return Is(Punct, rune)
}

// IsSpace reports whether the rune is a space character as defined
// by Unicode's White Space property; in the Latin-1 space
// this is
//	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
// Other definitions of spacing characters are set by category
// Z and property Pattern_White_Space.
func IsSpace(rune int) bool {
	// This property isn't the same as Z; special-case it.
	if uint32(rune) <= MaxLatin1 {
		switch rune {
		case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
			return true
		}
		return false
	}
	return Is(White_Space, rune)
}

// IsSymbol reports whether the rune is a symbolic character.
func IsSymbol(rune int) bool {
	if uint32(rune) <= MaxLatin1 {
		return properties[uint8(rune)]&pS != 0
	}
	return Is(Symbol, rune)
}
