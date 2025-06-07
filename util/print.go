package util

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

func PrintTable(rows [][]string, header []string) {
	table := tablewriter.NewWriter(os.Stdout)
	if len(header) > 0 {
		table.SetHeader(header)
	}
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}

func PrintJson(obj any) {
	c, _ := json.MarshalIndent(obj, "", "	")
	fmt.Println(string(c))
}
