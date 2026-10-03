package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) View() tea.View {
	return tea.NewView(baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n")
}
