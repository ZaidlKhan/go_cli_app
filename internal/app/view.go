package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) View() tea.View {
	if m.isEmpty {
		msg := lipgloss.NewStyle().
			Width(66).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("240")).
			Render("No active development servers found.\n\nPress 'q' to quit.")

		return tea.NewView(baseStyle.Render(msg) + "\n")
	}

	return tea.NewView(baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n")
}
