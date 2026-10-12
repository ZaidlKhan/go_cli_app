package app

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestTableNavigationKeysMoveSelection(t *testing.T) {
	m := NewModel()
	m.table.SetRows([]table.Row{
		{"first"},
		{"second"},
		{"third"},
	})

	if !m.table.Focused() {
		t.Fatal("table should be focused so it can receive navigation keys")
	}

	keys := []rune{tea.KeyDown, tea.KeyDown, tea.KeyUp}
	wantCursors := []int{1, 2, 1}
	for i, key := range keys {
		msg := tea.KeyPressMsg(tea.Key{Code: key})
		updated, _ := m.Update(msg)
		m = updated.(Model)
		if got := m.table.Cursor(); got != wantCursors[i] {
			t.Fatalf("cursor after key %q = %d, want %d", msg.String(), got, wantCursors[i])
		}
	}
}

func TestClearNoticeRestoresTableView(t *testing.T) {
	m := NewModel()
	m.setTableRows([]table.Row{{"3000", "dev-server", "1234", "http://localhost:3000", "main"}})
	m.state = StateKillSuccess

	updated, _ := m.Update(clearNoticeMsg{})
	m = updated.(Model)

	if m.state != StateNormal {
		t.Fatalf("state after clearing success notice = %v, want %v", m.state, StateNormal)
	}

	view := m.View()
	if !view.AltScreen {
		t.Fatal("app views should use the alternate screen to prevent stale content after redraws")
	}
	if !strings.Contains(view.Content, "Port") || !strings.Contains(view.Content, "dev-server") {
		t.Fatalf("table view not restored after notice: %q", view.Content)
	}
}

func TestReloadWithOneRemainingProcessShowsItsRow(t *testing.T) {
	m := NewModel()
	m.setTableRows([]table.Row{{"3001", "remaining-app", "1235", "http://localhost:3001", "main"}})

	if got, want := m.table.Height(), 1; got != want {
		t.Fatalf("table viewport height = %d, want %d", got, want)
	}
	if !strings.Contains(m.table.View(), "remaining-app") {
		t.Fatal("table should render the remaining process when only one row is left")
	}
}

func TestInitialTableHeightMatchesRows(t *testing.T) {
	m := NewModel()

	if got, want := m.table.Height(), len(m.table.Rows()); got != want {
		t.Fatalf("initial table viewport height = %d, want %d rows", got, want)
	}
}

func TestInitSchedulesProcessPolling(t *testing.T) {
	m := Model{}
	if m.Init() == nil {
		t.Fatal("Init should schedule the first process poll")
	}
}

func TestProcessPollUpdatesRowsAndSchedulesNextPoll(t *testing.T) {
	m := Model{table: table.New()}
	m.table.SetColumns([]table.Column{
		{Title: "Port", Width: 6},
		{Title: "App", Width: 14},
		{Title: "PID", Width: 6},
		{Title: "URL", Width: 24},
		{Title: "Branch", Width: 10},
	})
	m.setTableRows([]table.Row{{"3000", "first-app", "1234", "http://localhost:3000", "main"}})

	updated, cmd := m.Update(processRowsMsg{rows: []table.Row{
		{"3001", "second-app", "1235", "http://localhost:3001", "main"},
	}})
	m = updated.(Model)

	if got := m.table.Rows()[0][1]; got != "second-app" {
		t.Fatalf("poll result app name = %q, want %q", got, "second-app")
	}
	if cmd == nil {
		t.Fatal("processing poll results should schedule the next poll")
	}
}

func TestKillPromptAndSuccessMessageUseSameSizeBox(t *testing.T) {
	m := Model{state: StateConfirmKill, targetName: "dev-server", targetPort: "3000"}
	promptView := strings.SplitN(m.View().Content, "\n  ", 2)[0]

	m.state = StateKillSuccess
	m.targetPID = 1234
	successView := strings.SplitN(m.View().Content, "\n  ", 2)[0]

	if got, want := lipgloss.Width(successView), lipgloss.Width(promptView); got != want {
		t.Errorf("success box width = %d, want prompt box width %d", got, want)
	}
	if got, want := lipgloss.Height(successView), lipgloss.Height(promptView); got != want {
		t.Errorf("success box height = %d, want prompt box height %d", got, want)
	}
}
