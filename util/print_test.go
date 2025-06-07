package util

import "testing"

func TestPrintTable(t *testing.T) {
	var rows [][]string
	rows = append(rows, []string{"1", "2", "3"})
	rows = append(rows, []string{"1", "2", "3"})
	rows = append(rows, []string{"1", "2", "3"})

	PrintTable(rows, nil)
}
