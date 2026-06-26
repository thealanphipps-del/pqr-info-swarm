package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type StreamAppendMsg string

type StreamPaneModel struct {
	viewport viewport.Model
	content  string
}

func NewStreamPaneModel() StreamPaneModel {
	v := viewport.New(60, 12)
	v.SetContent("Live Swarm Stream\n\nWaiting for telemetry...")
	return StreamPaneModel{viewport: v}
}

func (m StreamPaneModel) Update(msg tea.Msg) (StreamPaneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case StreamAppendMsg:
		m.content += string(msg) + "\n\n"
		m.viewport.SetContent("Live Swarm Stream\n\n" + m.content)
		m.viewport.GotoBottom()
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m StreamPaneModel) View() string {
	return m.viewport.View()
}

func updateStream(m StreamPaneModel, msg tea.Msg, cmds []tea.Cmd) (StreamPaneModel, []tea.Cmd) {
	var cmd tea.Cmd
	m, cmd = m.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, cmds
}
