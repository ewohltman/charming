package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ewohltman/charming/pkg/list"
	"github.com/ewohltman/charming/pkg/viewport"
)

type Model struct {
	listModel     list.Model
	viewportModel viewport.Model
}

func NewModel() Model {
	return Model{
		listModel:     list.NewModel(),
		viewportModel: viewport.NewModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.listModel.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	listModel, listCmd := m.listModel.Update(msg)
	viewportModel, viewportCmd := m.viewportModel.Update(msg)

	m.listModel = listModel.(list.Model)
	m.viewportModel = viewportModel.(viewport.Model)

	return m, tea.Batch(
		listCmd,
		viewportCmd,
	)
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s\n",
		m.listModel.View(),
		m.viewportModel.View(),
	)
}
