package app

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func actionBox(content string, borderColor string) string {
	return lipgloss.NewStyle().
		Width(66).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(1, 2).
		Render(content)
}

func appView(content string) tea.View {
	view := tea.NewView(content)
	view.AltScreen = true
	return view
}

func (m Model) View() tea.View {
	const containerWidth = 66

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	var helpText string
	switch m.state {
	case StateConfirmKill:
		helpText = helpStyle.Render("y / ↵ : confirm  •  n: cancel")
	case StateKillSuccess, StateKillError, StateOpenError:
		helpText = helpStyle.Render("press any key to continue")
	case StateNormal:
		if m.isEmpty {
			helpText = helpStyle.Render("q: quit")
		} else {
			helpText = helpStyle.Render("k: kill  •  o: open  •  q: quit")
		}
	}

	if m.state == StateKillSuccess {
		msg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Render(fmt.Sprintf("✔ Successfully terminated %s (PID %d)", m.targetName, m.targetPID))

		return appView(actionBox(msg, "42") + "\n  " + helpText + "\n")
	}

	if m.state == StateKillError {
		msg := lipgloss.NewStyle().
			Width(containerWidth).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("204")).
			Render(fmt.Sprintf("❌ Failed to kill PID %d: %s", m.targetPID, m.errMessage))

		return appView(actionBox(msg, "204") + "\n  " + helpText + "\n")
	}

	if m.state == StateOpenError {
		msg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("204")).
			Render(fmt.Sprintf("❌ Failed to open URL: %s", m.errMessage))

		return appView(actionBox(msg, "204") + "\n  " + helpText + "\n")
	}

	if m.state == StateConfirmKill {
		boldStyle := lipgloss.NewStyle().Bold(true)
		appName := boldStyle.Render(m.targetName)
		port := boldStyle.Render(m.targetPort)

		content := fmt.Sprintf(
			"Kill %s running on port %s?",
			appName, port,
		)

		return appView(actionBox(content, "204") + "\n  " + helpText + "\n")
	}

	if m.isEmpty {
		msg := lipgloss.NewStyle().
			Width(containerWidth).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("240")).
			Render("No active development servers found.")

		return appView(baseStyle.Render(msg) + "\n  " + helpText + "\n")
	}

	return appView(baseStyle.Render(m.table.View()) + "\n  " + helpText + "\n")
}
