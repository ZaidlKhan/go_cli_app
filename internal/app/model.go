package app

import (
	"os/exec"
	"sort"
	"strconv"

	"main/internal/parser"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

func loadProcessRows() []table.Row {
	cmdStruct := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P", "-Fpcn")
	out, _ := cmdStruct.CombinedOutput()
	processes := parser.ParseOutput(string(out))

	rows := make([]table.Row, 0, len(processes))
	for _, proc := range processes {
		rows = append(rows, table.Row{
			strconv.Itoa(int(proc.Port)),
			proc.ProcessName,
			strconv.Itoa(proc.PID),
			proc.URL,
			proc.GitBranch,
		})
	}

	sortProcessRows(rows)
	return rows
}

func sortProcessRows(rows []table.Row) {
	sort.Slice(rows, func(i, j int) bool {
		portI, _ := strconv.Atoi(rows[i][0])
		portJ, _ := strconv.Atoi(rows[j][0])
		if portI != portJ {
			return portI < portJ
		}
		if rows[i][1] != rows[j][1] {
			return rows[i][1] < rows[j][1]
		}
		return rows[i][2] < rows[j][2]
	})
}

type State int

const (
	StateNormal State = iota
	StateConfirmKill
	StateKillSuccess
	StateKillError
	StateOpenError
)

type Model struct {
	table      table.Model
	isEmpty    bool
	state      State
	targetPID  int
	targetPort string
	targetName string
	errMessage string
}

func (m *Model) reloadTable() {
	columns := []table.Column{
		{Title: "Port", Width: 6},
		{Title: "App", Width: 14},
		{Title: "PID", Width: 6},
		{Title: "URL", Width: 24},
		{Title: "Branch", Width: 10},
	}

	m.table.SetColumns(columns)
	m.setTableRows(loadProcessRows())
}

func (m *Model) setTableRows(rows []table.Row) {
	m.isEmpty = len(rows) == 0
	if m.isEmpty {
		rows = []table.Row{{"", "", "", "", ""}}
	}

	m.table.SetRows(rows)
	m.table.SetWidth(70)
	m.table.SetHeight(len(rows) + 2)
}

func NewModel() Model {
	m := Model{table: table.New()}

	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(lipgloss.Color("#C4B5FD")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#8B5CF6")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#F3E8FF")).
		Background(lipgloss.Color("#4C1D95")).
		Bold(false)
	m.table.SetStyles(s)
	m.reloadTable()
	m.table.Focus()

	return m
}
