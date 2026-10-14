package app

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableNavigationKeysMoveSelection(t *testing.T) {
	m := NewModel()
	m.table.SetRows([]table.Row{
		{"first"},
		{"second"},
		{"third"},
	})

	require.True(t, m.table.Focused(), "table should be focused so it can receive navigation keys")

	keys := []rune{tea.KeyDown, tea.KeyDown, tea.KeyUp}
	wantCursors := []int{1, 2, 1}
	for i, key := range keys {
		msg := tea.KeyPressMsg(tea.Key{Code: key})
		updated, _ := m.Update(msg)
		m = updated.(Model)
		assert.Equal(t, wantCursors[i], m.table.Cursor(), "cursor after key %q", msg.String())
	}
}

func TestClearNoticeRestoresTableView(t *testing.T) {
	m := NewModel()
	m.setTableRows([]table.Row{{"3000", "dev-server", "1234", "http://localhost:3000", "main"}})
	m.state = StateKillSuccess

	updated, _ := m.Update(clearNoticeMsg{})
	m = updated.(Model)

	assert.Equal(t, StateNormal, m.state, "state after clearing success notice")

	view := m.View()
	assert.True(t, view.AltScreen, "app views should use the alternate screen to prevent stale content after redraws")
	assert.Contains(t, view.Content, "Port")
	assert.Contains(t, view.Content, "dev-server")
}

func TestReloadWithOneRemainingProcessShowsItsRow(t *testing.T) {
	m := NewModel()
	m.setTableRows([]table.Row{{"3001", "remaining-app", "1235", "http://localhost:3001", "main"}})

	assert.Equal(t, 1, m.table.Height(), "table viewport height")
	assert.Contains(t, m.table.View(), "remaining-app", "table should render the remaining process when only one row is left")
}

func TestInitialTableHeightMatchesRows(t *testing.T) {
	m := NewModel()

	assert.Equal(t, len(m.table.Rows()), m.table.Height(), "initial table viewport height matches row count")
}

func TestInitSchedulesProcessPolling(t *testing.T) {
	m := Model{}
	assert.NotNil(t, m.Init(), "Init should schedule the first process poll")
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

	assert.Equal(t, "second-app", m.table.Rows()[0][1], "poll result app name")
	assert.NotNil(t, cmd, "processing poll results should schedule the next poll")
}

func TestSortProcessRowsKeepsPollingOrderStable(t *testing.T) {
	rows := []table.Row{
		{"9000", "third-app", "300", "http://localhost:9000", "main"},
		{"8080", "second-app", "200", "http://localhost:8080", "main"},
		{"3000", "first-app", "100", "http://localhost:3000", "main"},
	}

	sortProcessRows(rows)

	assert.Equal(t, []table.Row{
		{"3000", "first-app", "100", "http://localhost:3000", "main"},
		{"8080", "second-app", "200", "http://localhost:8080", "main"},
		{"9000", "third-app", "300", "http://localhost:9000", "main"},
	}, rows)
}

func TestKillPromptAndSuccessMessageUseSameSizeBox(t *testing.T) {
	m := Model{state: StateConfirmKill, targetName: "dev-server", targetPort: "3000"}
	promptView := strings.SplitN(m.View().Content, "\n  ", 2)[0]

	m.state = StateKillSuccess
	m.targetPID = 1234
	successView := strings.SplitN(m.View().Content, "\n  ", 2)[0]

	assert.Equal(t, lipgloss.Width(promptView), lipgloss.Width(successView), "success box width should match prompt box width")
	assert.Equal(t, lipgloss.Height(promptView), lipgloss.Height(successView), "success box height should match prompt box height")
}

func TestOpenURLFailureIsShownToUser(t *testing.T) {
	m := Model{table: table.New()}
	m.table.SetColumns([]table.Column{
		{Title: "Port", Width: 6},
		{Title: "App", Width: 14},
		{Title: "PID", Width: 6},
		{Title: "URL", Width: 24},
		{Title: "Branch", Width: 10},
	})
	m.setTableRows([]table.Row{{"3000", "dev-server", "1234", "ftp://localhost:3000", "main"}})

	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Text: "o", Code: 'o'}))
	m = updated.(Model)

	assert.Equal(t, StateOpenError, m.state, "state after unsupported URL")
	assert.Contains(t, m.View().Content, "Failed to open URL", "open URL error should be shown in the app view")
}
