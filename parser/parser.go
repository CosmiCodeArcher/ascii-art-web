package parser

import (
	"fmt"
	"strings"
)

func ParseInput(text string) ([]string, error) {
	if text == "" {
		return nil, fmt.Errorf("Empty Input")
	}
	return strings.Split(strings.ReplaceAll(text, "\\n", "\n"), "\n"), nil
}