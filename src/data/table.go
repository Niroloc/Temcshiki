package data

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
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

func (this *Table) Render() string {
	tw := table.NewWriter()
	titles := table.Row{}
	for _, t := range this.Titles {
		titles = append(titles, t)
	}
	tw.AppendHeader(titles)
	for _, r := range this.Values {
		row := table.Row{}
		for _, e := range r {
			row = append(row, e)
		}
		tw.AppendRow(row)
	}
	return tw.Render()
}
