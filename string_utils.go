package ddderr

import (
	"strings"
	"unicode"
)

func getSanitizedStatusName(attribute, operation string) string {
	if attribute == "" {
		return operation
	}
	return replaceUnderscoreFromString(attribute) + operation
}

func replaceUnderscoreFromString(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	isTitle := false
	for i, ch := range str {
		if unicode.IsLetter(ch) {
			if i == 0 || isTitle {
				isTitle = false
				ch = unicode.ToTitle(ch)
			}
			b.WriteRune(ch)
			continue
		}
		isTitle = true
	}
	return b.String()
}
