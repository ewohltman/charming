package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
)

type Row struct {
	Object    string
	Status    string
	RunStatus string
	Age       string
	Warnings  string
}

func (r Row) String() string {
	return fmt.Sprintf("Object: %s\nAge: %s",
		r.Object,
		r.Age,
	)
}

func (r Row) convertToTableRow() table.Row {
	return table.Row{
		r.Object,
		r.Status,
		r.RunStatus,
		r.Age,
		r.Warnings,
	}
}
