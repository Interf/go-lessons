package pack

import (
	"strconv"
)

func Pack(text string) string {

	var result string

	if len(text) == 0 {
		return result
	}

	countMap := make(map[rune]int)

	for _, char := range text {
		countMap[char]++
	}

	var prev rune

	for i := 0; i < len(text); i++ {

		if prev != rune(text[i]) {
			result += string(text[i])

			if countMap[rune(text[i])] > 1 {
				result += strconv.Itoa(countMap[rune(text[i])])
			}

			prev = rune(text[i])
		}

	}

	return result
}
