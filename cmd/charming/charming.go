package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ewohltman/charming/pkg/model"
)

func main() {
	items := []list.Item{
		model.Item("Ramen"),
		model.Item("Tomato Soup"),
		model.Item("Hamburgers"),
		model.Item("Cheeseburgers"),
		model.Item("Currywurst"),
		model.Item("Okonomiyaki"),
		model.Item("Pasta"),
		model.Item("Fillet Mignon"),
		model.Item("Caviar"),
		model.Item("Just Wine"),
	}

	const defaultWidth = 20

	l := list.New(items, model.ItemDelegate{}, defaultWidth, model.ListHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = model.TitleStyle
	l.Styles.PaginationStyle = model.PaginationStyle
	l.Styles.HelpStyle = model.HelpStyle

	m := model.Model{List: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
