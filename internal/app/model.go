package app

import (
	"os/exec"
	"strconv"

	"main/internal/parser"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

type Model struct {
	table table.Model
}

func NewModel() Model {
	columns := []table.Column{
		{Title: "Port", Width: 6},
		{Title: "App", Width: 14},
		{Title: "PID", Width: 6},
		{Title: "URL", Width: 24},
		{Title: "Branch", Width: 14},
	}

	cmdStruct := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P", "-Fpcn")
	out, _ := cmdStruct.CombinedOutput()
	processes := parser.ParseOutput(string(out))

	var rows []table.Row
	for _, proc := range processes {
		rows = append(rows, table.Row{
			strconv.Itoa(int(proc.Port)),
			proc.ProcessName,
			strconv.Itoa(proc.PID),
			proc.URL,
			proc.GitBranch,
		})
	}

	if len(rows) == 0 {
		rows = append(rows, table.Row{"-", "No dev servers found", "-", "-", "-"})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
		table.WithWidth(120),
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
