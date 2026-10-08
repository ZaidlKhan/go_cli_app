package app

import (
	"os/exec"

	"main/internal/filesystems"

	"charm.land/bubbles/v2/table"
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
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
	}

	cmdStruct := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P", "-Fpcn")
	out, _ := cmdStruct.CombinedOutput()
	filesystems.ParseOutput(string(out))

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

	return Model{table: t}
}
