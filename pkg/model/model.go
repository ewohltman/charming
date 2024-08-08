package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	margin  = 2
	padding = 4
	width   = 20

	listHeight     = 24
	viewportHeight = 2
)

var (
	titleStyle      = list.DefaultStyles().Title
	paginationStyle = list.DefaultStyles().PaginationStyle.
			PaddingLeft(padding).
			PaddingRight(padding)
	helpStyle = list.DefaultStyles().HelpStyle.
			PaddingLeft(padding).
			PaddingRight(padding)
)

type Client interface {
	List() ([]Item, error)
	Get(name string) (Item, error)
}

type UpdateListMsg struct {
	Items []Item
}

func (ulm UpdateListMsg) convertToListItems() []list.Item {
	converted := make([]list.Item, len(ulm.Items))

	for i := range ulm.Items {
		converted[i] = ulm.Items[i]
	}

	return converted
}

type Model struct {
	client        Client
	listModel     list.Model
	viewportModel viewport.Model
	quitting      bool
}

func New(client Client) *Model {
	listModel := list.New(nil, ItemDelegate{}, width, listHeight)
	listModel.Title = "What do you want for dinner?"
	listModel.Styles.Title = titleStyle
	listModel.Styles.PaginationStyle = paginationStyle
	listModel.Styles.HelpStyle = helpStyle

	listModel.SetShowStatusBar(false)
	listModel.SetFilteringEnabled(false)

	viewportModel := viewport.New(width, viewportHeight)
	viewportModel.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingLeft(padding).
		PaddingRight(padding)

	return &Model{
		client:        client,
		listModel:     listModel,
		viewportModel: viewportModel,
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
		m.listModel.View(),
		m.viewportModel.View(),
	)

	if m.quitting {
		view += "\n"
	}

	return view
}

func (m *Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (cmd tea.Cmd) {
	m.listModel.SetWidth(msg.Width - margin)
	m.listModel.SetHeight(msg.Height - 6)
	m.viewportModel.Width = msg.Width - margin

	return m.updateModels(msg)
}

func (m *Model) handleUpdateListMsg(msg UpdateListMsg) (cmd tea.Cmd) {
	if cmd = m.listModel.SetItems(msg.convertToListItems()); cmd != nil {
		m.updateModels(cmd())
	}

	return nil
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (cmd tea.Cmd) {
	switch key := msg.String(); key {
	case "q", "ctrl+c", "esc":
		m.quitting = true

		return tea.Quit
	default:
		item, err := m.client.Get(string(m.listModel.SelectedItem().(Item)))
		if err != nil {
			panic("error getting item: " + err.Error())
		}

		m.viewportModel.SetContent(fmt.Sprintf("%s", item))

		return m.updateModels(msg)
	}
}

func (m *Model) updateModels(msg tea.Msg) (cmd tea.Cmd) {
	var cmds []tea.Cmd

	if m.listModel, cmd = m.listModel.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	if m.viewportModel, cmd = m.viewportModel.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tea.Sequence(cmds...)
}
