package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/erik-adelbert/firework/internal/sym"
)

func (m *Model) Init() tea.Cmd {
	m.last = time.Now()

	return tick()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		now := time.Time(msg)

		m.step(now)
		m.sample(now) // for FPS calculation

		return m, tick()

	case tea.WindowSizeMsg:
		h, w := msg.Height, msg.Width
		sz := h * w

		if sz > cap(m.screen) {
			m.screen = make([]sym.Symbol, sz)
		}
		m.screen = m.screen[:sz]
		m.h, m.w = h, w

	case tea.KeyMsg:
		switch m.mode {
		case modeDemo:
			return m.demoKeys(msg)
		default:
			return m.defaultKeys(msg)
		}
	}

	return m, nil
}

func (m *Model) demoKeys(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) defaultKeys(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "m":
		m.cycleMode()
	case "p":
		m.togglePause()
	case "r":
		m.resetNow()
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(30*ms, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

const ms = time.Millisecond
