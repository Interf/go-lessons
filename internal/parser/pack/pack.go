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
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {

		if prev != runes[i] {
			result += string(runes[i])

			if countMap[runes[i]] > 1 {
				result += strconv.Itoa(countMap[runes[i]])
			}

			prev = runes[i]
		}

	}

	return result
}
