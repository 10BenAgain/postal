package utils

import (
	"fmt"
	"strings"
	"unicode"
)

var (
	ErrOutOfRange          = fmt.Errorf("input value out of range")
	ErrPKIncorrectFileSize = fmt.Errorf("pk* file does not match expected value")
)

type UNumber interface {
	uint | uint8 | uint16 | uint32 | uint64
}

func GetSliceFromRawData(s []byte, offset int, length int) []byte {
	if offset+length > len(s) {
		return nil
	}
	return s[offset : offset+length]
}

func WordValuesToPID(word uint16) uint32 {
	return (uint32(word) << 16) | uint32(word)
}

func Bool2Byte(b bool) byte {
	if b {
		return 1
	}

	return 0
}

func Uint2Bool[N UNumber](u N) bool {
	return u > 0
}

func PrintHexFromByteSlice(bs []byte) {
	for i, v := range bs {
		if i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("0x%02X, ", v)
	}
	fmt.Println()
}

func SanitizeSearch(s string) string {
	if len(s) == 0 {
		return ""
	}

	ret := new(strings.Builder)
	split := strings.Split(s, " ")
	single := false

	if len(split) == 1 {
		split = strings.Split(s, "-")
		if len(split) > 1 {
			single = true
		}
	}

	for i, word := range split {
		r := []rune(word)
		for j := range word {
			if j == 0 {
				r[0] = unicode.ToUpper(r[0])
			} else {
				r[j] = unicode.ToLower(r[j])
			}
		}

		if i < len(split)-1 && !single {
			r = append(r, ' ')
		} else if i < len(split)-1 && single {
			r = append(r, '-')
		}

		ret.WriteString(string(r))
	}

	return ret.String()
}

func SanitizeWordSearch(s string) string {
	if len(s) == 0 {
		return ""
	}

	r := []rune(s)
	for i := range r {
		r[i] = unicode.ToUpper(r[i])
	}

	return string(r)
}
