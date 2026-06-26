package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatAppendMsg string

type ChatPaneModel struct {
	history viewport.Model
	input   textarea.Model
	content string
}

func NewChatPaneModel() ChatPaneModel {
	h := viewport.New(60, 10)
	h.SetContent("Agent Chat\n\nWaiting for conversation...")

	in := textarea.New()
	in.Placeholder = "Type message to agent..."
	in.SetHeight(3)

	return ChatPaneModel{
		history: h,
		input:   in,
	}
}

func (m ChatPaneModel) Update(msg tea.Msg) (ChatPaneModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ChatAppendMsg:
		m.content += string(msg) + "\n\n"
		m.history.SetContent("Agent Chat\n\n" + m.content)
		m.history.GotoBottom()
		return m, nil
	}

	var cmd tea.Cmd
	m.history, cmd = m.history.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.input, cmd = m.input.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m ChatPaneModel) View() string {
	return m.history.View() + "\n" + m.input.View()
}

func updateChat(m ChatPaneModel, msg tea.Msg, cmds []tea.Cmd) (ChatPaneModel, []tea.Cmd) {
	var cmd tea.Cmd
	m, cmd = m.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, cmds
}
