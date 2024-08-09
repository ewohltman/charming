package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ewohltman/charming/pkg/client"
	"github.com/ewohltman/charming/pkg/model"
)

func main() {
	c := client.NewFake()
	m := model.New(c)
	p := tea.NewProgram(m)

	go c.Watch(p)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
