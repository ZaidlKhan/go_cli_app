package app

import (
	"strconv"
	"time"

	"main/internal/process"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

const processPollInterval = 2500 * time.Millisecond

type clearNoticeMsg struct{}
type processPollMsg struct{}
type processRowsMsg struct {
	rows []table.Row
}

func (m Model) Init() tea.Cmd {
	return scheduleProcessPoll()
}

func scheduleProcessPoll() tea.Cmd {
	return tea.Tick(processPollInterval, func(time.Time) tea.Msg {
		return processPollMsg{}
	})
}

func pollProcessRows() tea.Cmd {
	return func() tea.Msg {
		return processRowsMsg{rows: loadProcessRows()}
	}
}

func rowsEqual(a, b []table.Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case processPollMsg:
		return m, pollProcessRows()

	case processRowsMsg:
		rows := msg.rows
		if len(rows) == 0 {
			rows = []table.Row{{"", "", "", "", ""}}
		}
		if m.isEmpty != (len(msg.rows) == 0) || !rowsEqual(m.table.Rows(), rows) {
			m.setTableRows(msg.rows)
		}
		return m, scheduleProcessPoll()

	case clearNoticeMsg:
		m.state = StateNormal
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case StateKillSuccess, StateKillError:
			m.state = StateNormal
			return m, nil

		case StateNormal:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "k":
				if m.isEmpty {
					return m, nil
				}

				selectedRow := m.table.SelectedRow()
				if len(selectedRow) >= 3 {
					m.targetPort = selectedRow[0]
					m.targetName = selectedRow[1]

					pid, err := strconv.Atoi(selectedRow[2])
					if err == nil {
						m.targetPID = pid
						m.state = StateConfirmKill
					}
				}
				return m, nil
			}

		case StateConfirmKill:
			switch msg.String() {
			case "y", "Y", "enter":
				err := process.KillProcess(m.targetPID)
				if err != nil {
					m.errMessage = err.Error()
					m.state = StateKillError
				} else {
					m.reloadTable()
					m.state = StateKillSuccess
					return m, tea.Tick(1000*time.Millisecond, func(t time.Time) tea.Msg {
						return clearNoticeMsg{}
					})
				}
				return m, nil

			case "n", "N", "esc":
				m.state = StateNormal
				return m, nil
			}
		}
	}

	if m.state == StateNormal {
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}
