package main_test

import (
	"strings"
)

func parseTable(out string) [][]string {
	var parsed [][]string
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") {
			continue
		}
		if strings.HasPrefix(line, "|") {
			cells := strings.Split(line, "|")
			var row []string
			for i := 1; i < len(cells)-1; i++ {
				row = append(row, strings.TrimSpace(cells[i]))
			}
			parsed = append(parsed, row)
		}
	}
	return parsed
}
