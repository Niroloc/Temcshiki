package data

import (
	"fmt"
)

type Table struct {
	Titles []string
	Values [][]string
}

func CreateNewTable(titles []string) Table {
	return Table{titles, [][]string{}}
}

func (this *Table) AddLine(line []string) error {
	if len(line) != len(this.Titles) {
		return fmt.Errorf("Columns number mismatch: expected %d, given %d", len(this.Titles), len(line))
	}
	this.Values = append(this.Values, line)
	return nil
}
