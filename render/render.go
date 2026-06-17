package render

import "strings"

func RenderToString(charMap map[rune][]string, lines []string) string {
	var result strings.Builder

	for _, line := range lines {

		if line == "" {
			result.WriteString("\n")
			continue
		}

		for row := 0; row < 8; row++ {
			var rowBuilder strings.Builder

			for _, char := range line {
				if charRowSlice, ok := charMap[char]; ok && len(charRowSlice) > row {
					rowBuilder.WriteString(charRowSlice[row])
				}
			}
			result.WriteString(rowBuilder.String())
			result.WriteByte('\n')
		}
	}

	return result.String()
}