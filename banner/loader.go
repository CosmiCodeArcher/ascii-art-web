package banner

import (
	"io"
	"os"
	"strings"
)

func LoadBanner(filename string) (map[rune][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fileContentStr := string(fileContent)
	fileContentStr = strings.ReplaceAll(fileContentStr, "\r\n", "\n")
	chunks := strings.Split(strings.Trim(fileContentStr, "\n"), "\n\n")

	charMap := make(map[rune][]string)
	for i, chunk := range chunks {
		r := rune(32 + i)
		charMap[r] = strings.Split(strings.Trim(chunk, "\n"), "\n")
	}

	return charMap, nil
}