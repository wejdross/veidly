package helpers

import (
	"strings"
	"unicode/utf8"
)

func NormLevenshtein(a, b string) float32 {
	if len(b) == 0 || len(a) == 0 {
		return 0
	}
	l := float32(Levenshtein(a, b))
	var maxLen float32 = 0
	if len(a) > len(b) {
		maxLen = float32(len(a))
	} else {
		maxLen = float32(len(b))
	}
	return 1 - (l / maxLen)
}

func min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func Levenshtein(a, b string) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	f := make([]int, utf8.RuneCountInString(b)+1)
	for j := range f {
		f[j] = j
	}
	for _, ca := range a {
		j := 1
		fj1 := f[0]
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1)
			if cb != ca {
				mn = min(mn, fj1+1)
			} else {
				mn = min(mn, fj1)
			}
			fj1, f[j] = f[j], mn
			j++
		}
	}
	return f[len(f)-1]
}
