package client

import (
	"errors"
	"math/rand/v2"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ewohltman/charming/pkg/model"
)

type Fake struct {
	rows   []model.Row
	rowMap map[string]model.Row
}

func NewFake() *Fake {
	rows := []model.Row{
		{Object: "Object 1", Status: "(running)", RunStatus: "planning", Age: "", Warnings: ""},
		{Object: "Object 2", Status: "READY", RunStatus: "planned_and_finished", Age: "", Warnings: ""},
		{Object: "Object 3", Status: "READY", RunStatus: "planned_and_finished", Age: "", Warnings: ""},
		{Object: "Object 4", Status: "READY", RunStatus: "planned_and_finished", Age: "", Warnings: ""},
		{Object: "Object 5", Status: "READY", RunStatus: "planned_and_finished", Age: "", Warnings: ""},
		{Object: "Object 6", Status: "READY", RunStatus: "planned_and_finished", Age: "", Warnings: ""},
		{Object: "Object 7", Status: "(running)", RunStatus: "planning", Age: "", Warnings: ""},
		{Object: "Object 8", Status: "READY", RunStatus: "planned_and_finished", Age: "", Warnings: ""},
		{Object: "Object 9", Status: "(running)", RunStatus: "planning", Age: "", Warnings: ""},
		{Object: "Object 10", Status: "(running)", RunStatus: "planning", Age: "", Warnings: ""},
	}

	rowMap := make(map[string]model.Row)

	for _, row := range rows {
		rowMap[row.Object] = row
	}

	return &Fake{
		rows:   rows,
		rowMap: rowMap,
	}
}

func (f *Fake) List() ([]model.Row, error) {
	for i := range f.rows {
		f.rows[i].Age = time.Duration(rand.Int64N(90000000000)).String()
		f.rowMap[f.rows[i].Object] = f.rows[i]
	}

	return f.rows, nil
}

func (f *Fake) Get(name string) (model.Row, error) {
	name = strings.TrimSpace(name)

	row, found := f.rowMap[name]
	if !found {
		return model.Row{}, errors.New("row not found")
	}

	return row, nil
}

func (f *Fake) Watch(program *tea.Program) {
	items, _ := f.List()

	program.Send(model.UpdateListMsg{Items: items})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		items, _ = f.List()

		program.Send(model.UpdateListMsg{Items: items})
	}
}
