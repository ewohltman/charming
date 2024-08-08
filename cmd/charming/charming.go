package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ewohltman/charming/pkg/client"
	"github.com/ewohltman/charming/pkg/model"
)

func watch(p *tea.Program, c *client.Fake) {
	items, _ := c.List()

	p.Send(model.UpdateListMsg{Items: items})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		items, _ = c.List()

		p.Send(model.UpdateListMsg{Items: items})
	}
}

func main() {
	c := client.NewFake()
	m := model.New(c)
	p := tea.NewProgram(m)

	go watch(p, c)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
