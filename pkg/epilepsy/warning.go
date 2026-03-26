// Package epilepsy provides a warning message for users with photosensitive
// epilepsy.
package epilepsy

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var Warning = []string{
	"⚠️ WARNING: This program displays rapidly flashing images.",
	"It may potentially trigger seizures for people with photosensitive epilepsy.",
	"User discretion is advised.",
	"Proceed [y]/N?",
}

type model struct {
	confirmed bool
}

func NewModel() *model {
	return &model{}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.confirmed {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirmed = true
		}

		return m, tea.Quit
	}

	return m, nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) View() tea.View {
	var sb strings.Builder

	style := lipgloss.NewStyle().Padding(1, 2).Bold(true)

	sb.WriteString(style.Render(strings.Join(Warning, "\n")))

	return tea.NewView(sb.String())
}

func Warn() bool {
	m := NewModel()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running warning program:", err)
		return false
	}

	return m.confirmed
}
