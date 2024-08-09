package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	margin  = 2
	padding = 4

	viewportHeight = 2
)

type Client interface {
	List() ([]Row, error)
	Get(name string) (Row, error)
}

type UpdateListMsg struct {
	Items []Row
}

func (ulm UpdateListMsg) convertToTableRows() []table.Row {
	converted := make([]table.Row, len(ulm.Items))

	for i := range ulm.Items {
		converted[i] = ulm.Items[i].convertToTableRow()
	}

	return converted
}

type Model struct {
	client         Client
	tableBaseStyle lipgloss.Style
	tableModel     table.Model
	viewportModel  viewport.Model
	quitting       bool
}

func New(client Client) *Model {
	columns := []table.Column{
		{Title: "Object", Width: 0},
		{Title: "Status", Width: 0},
		{Title: "Run Status", Width: 0},
		{Title: "Age", Width: 0},
		{Title: "Warnings", Width: 0},
	}

	tableBaseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tableStyle.Selected = tableStyle.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	tableModel.SetStyles(tableStyle)

	viewportModel := viewport.New(0, viewportHeight)
	viewportModel.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingLeft(padding).
		PaddingRight(padding)

	return &Model{
		client:         client,
		tableBaseStyle: tableBaseStyle,
		tableModel:     tableModel,
		viewportModel:  viewportModel,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Sequence(
		m.viewportModel.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		return m, m.handleWindowSizeMsg(msg)
	case UpdateListMsg:
		return m, m.handleUpdateListMsg(msg)
	default:
		return m, m.updateModels(msg)
	}
}

func (m *Model) View() string {
	view := fmt.Sprintf("%s\n%s",
		m.tableBaseStyle.Render(m.tableModel.View()),
		m.viewportModel.View(),
	)

	if m.quitting {
		view += "\n"
	}

	return view
}

func (m *Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (cmd tea.Cmd) {
	width := msg.Width - margin
	columnWidth := (width / 5) - margin

	columns := []table.Column{
		{Title: "Object", Width: columnWidth},
		{Title: "Status", Width: columnWidth},
		{Title: "Run Status", Width: columnWidth},
		{Title: "Age", Width: columnWidth},
		{Title: "Warnings", Width: columnWidth},
	}

	m.tableModel.SetWidth(width)
	m.tableModel.SetHeight(msg.Height - 9)
	m.tableModel.SetColumns(columns)

	m.viewportModel.Width = width

	return m.updateModels(msg)
}

func (m *Model) handleUpdateListMsg(msg UpdateListMsg) (cmd tea.Cmd) {
	m.tableModel.SetRows(msg.convertToTableRows())

	return nil
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (cmd tea.Cmd) {
	switch key := msg.String(); key {
	case "q", "ctrl+c", "esc":
		m.quitting = true

		return tea.Quit
	default:
		return m.updateModels(msg)
	}
}

func (m *Model) updateModels(msg tea.Msg) (cmd tea.Cmd) {
	var cmds []tea.Cmd

	if m.tableModel, cmd = m.tableModel.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	row, err := m.client.Get(m.tableModel.SelectedRow()[0])
	if err != nil {
		m.viewportModel.SetContent("Error: " + err.Error())
	} else {
		m.viewportModel.SetContent(fmt.Sprintf("%s", row))
	}

	if m.viewportModel, cmd = m.viewportModel.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tea.Sequence(cmds...)
}
