package strhelper

import (
	"strings"
	"unicode/utf8"
)

// Clean returns a string without control characters and without invalid unicode runes
func Clean(s string) string {
	// remove control chars
	b := make([]byte, len(s))
	var bl int
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 32 && c != 127 {
			b[bl] = c
			bl++
		}
	}

	str := string(b[:bl])

	// remove invalid runes
	if !utf8.ValidString(str) {
		v := make([]rune, 0, len(str))
		for i, r := range str {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(str[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		str = string(v)
	}
	str = strings.Trim(str, " ")
	return str
}

// Trim "0x" from a string. Do nothing if not present.
func Trim0x(PrefixedString string) string {
	return strings.TrimPrefix(PrefixedString, "0x")
}
