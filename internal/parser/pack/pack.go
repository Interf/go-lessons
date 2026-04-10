package pack

import (
	"strconv"
	"strings"
)

func Pack(text string) string {

	builder := strings.Builder{}

	if len(text) == 0 {
		return builder.String()
	}

	countMap := make(map[rune]int)

	for _, char := range text {
		countMap[char]++
	}

	var prev rune
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {

		if prev != runes[i] {
			builder.WriteString(string(runes[i]))

			if countMap[runes[i]] > 1 {
				builder.WriteString(strconv.Itoa(countMap[runes[i]]))
			}

			prev = runes[i]
		}

	}

	return builder.String()
}
