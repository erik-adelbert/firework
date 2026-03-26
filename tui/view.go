package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

var headers = []string{
	"FPS: %.1f [m]ode: %s  State: %s [p]ause  [r]eset  [q]uit",
	"FPS: %.1f [q]uit",
}

func (m *Model) View() tea.View {
	status := "running"
	if m.paused {
		status = "paused"
	}

	hud := fmt.Sprintf(headers[0], m.FPS(), m.mode, status)
	if m.mode == modeDemo {
		hud = fmt.Sprintf(headers[1], m.FPS())
	}

	v := tea.NewView(hud + "\n" + m.Render())
	v.AltScreen = true

	return v
}
