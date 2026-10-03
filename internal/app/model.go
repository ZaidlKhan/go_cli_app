package app

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	table table.Model
}

func NewModel() Model {
	columns := []table.Column{
		{Title: "Port", Width: 6},
		{Title: "PID", Width: 5},
		{Title: "URL", Width: 15},
		{Title: "App", Width: 10},
	}

	rows := []table.Row{
		{"8080", "86527", "localhost:8080", "python"},
		{"3000", "83618", "localhost:3000", "node"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(4),
		table.WithWidth(43),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("19")).
		Bold(false)
	t.SetStyles(s)

	m := Model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return Model{t}
}
