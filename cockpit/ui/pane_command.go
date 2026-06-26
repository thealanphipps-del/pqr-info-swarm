package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/pqr-info/pqr-info-swarm/cockpit/internal/adapter"
)

type CommandPaneModel struct {
	input   textinput.Model
	output  string
}

func NewCommandPaneModel() CommandPaneModel {
	in := textinput.New()
	in.Placeholder = "/help, /ticket, /swarm kill..."
	in.Focus()

	return CommandPaneModel{
		input:  in,
		output: "Command Palette\n\n(placeholder)",
	}
}

type CommandSubmitMsg string

func (m CommandPaneModel) Update(msg tea.Msg) (CommandPaneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			val := m.input.Value()
			m.input.SetValue("")
			return m, func() tea.Msg {
				return CommandSubmitMsg(val)
			}
		}
	case CommandSubmitMsg:
		// UI-only; actual execution wired via adapter in cockpit model
		m.output = "Last command: " + string(msg)
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m CommandPaneModel) View() string {
	return m.output + "\n\n" + m.input.View()
}

func updateCommand(m CommandPaneModel, msg tea.Msg, cmds []tea.Cmd, client *adapter.SwenClient) (CommandPaneModel, []tea.Cmd) {
	switch msg.(type) {
	case CommandSubmitMsg:
		// later: call client.ExecuteCommand(...)
	}
	var cmd tea.Cmd
	m, cmd = m.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, cmds
}
