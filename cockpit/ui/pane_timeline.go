package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type TimelineItem string

func (t TimelineItem) Title() string       { return string(t) }
func (t TimelineItem) Description() string { return "" }
func (t TimelineItem) FilterValue() string { return string(t) }

type TimelineAppendMsg string

type TimelinePaneModel struct {
	list list.Model
}

func NewTimelinePaneModel() TimelinePaneModel {
	items := []list.Item{}
	l := list.New(items, list.NewDefaultDelegate(), 40, 10)
	l.Title = "Timeline"
	return TimelinePaneModel{list: l}
}

func (m TimelinePaneModel) Update(msg tea.Msg) (TimelinePaneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case TimelineAppendMsg:
		m.list.InsertItem(len(m.list.Items()), TimelineItem(string(msg)))
		return m, nil
	}
	
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m TimelinePaneModel) View() string {
	return m.list.View()
}

func updateTimeline(m TimelinePaneModel, msg tea.Msg, cmds []tea.Cmd) (TimelinePaneModel, []tea.Cmd) {
	var cmd tea.Cmd
	m, cmd = m.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, cmds
}
