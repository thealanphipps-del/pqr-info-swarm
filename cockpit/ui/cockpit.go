package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"

	"github.com/pqr-info/pqr-info-swarm/cockpit/internal/adapter"
	"github.com/pqr-info/pqr-info-swarm/cockpit/internal/config"
)

type WSConnMsg *websocket.Conn
type WSDataMsg adapter.ProxyEvent
type WSErrMsg error

func connectWS(url string) tea.Cmd {
	return func() tea.Msg {
		c, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			return WSErrMsg(err)
		}
		return WSConnMsg(c)
	}
}

func readWS(c *websocket.Conn) tea.Cmd {
	return func() tea.Msg {
		var event adapter.ProxyEvent
		err := c.ReadJSON(&event)
		if err != nil {
			return WSErrMsg(err)
		}
		return WSDataMsg(event)
	}
}

type CockpitModel struct {
	cfg    config.TenantConfig
	styles Styles
	client *adapter.SwenClient
	wsConn *websocket.Conn

	stream   StreamPaneModel
	timeline TimelinePaneModel
	chat     ChatPaneModel
	command  CommandPaneModel
	telemetry TelemetryPaneModel
}

func NewCockpitModel(cfg config.TenantConfig) CockpitModel {
	styles := NewStyles(cfg)
	client := adapter.NewSwenClient(cfg.Backend.SwarmAPIURL, cfg.Backend.TicketAPIURL)

	return CockpitModel{
		cfg:      cfg,
		styles:   styles,
		client:   client,
		stream:   NewStreamPaneModel(),
		timeline: NewTimelinePaneModel(),
		chat:     NewChatPaneModel(),
		command:  NewCommandPaneModel(),
		telemetry: NewTelemetryPaneModel(),
	}
}

func (m CockpitModel) Init() tea.Cmd {
	// connect to proxy websocket
	return connectWS("ws://127.0.0.1:8081/ws")
}

func (m CockpitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case WSConnMsg:
		m.wsConn = msg
		m.timeline, _ = updateTimeline(m.timeline, TimelineAppendMsg("Connected to Swarm Proxy WS"), nil)
		cmds = append(cmds, readWS(m.wsConn))
	case WSDataMsg:
		// Dump for debugging
		_ = os.WriteFile("debug_payload.json", []byte(msg.RespBody), 0644)

		var completion adapter.ChatCompletion
		if err := json.Unmarshal([]byte(msg.RespBody), &completion); err == nil {
			if len(completion.Choices) > 0 {
				choice := completion.Choices[0].Message
				if choice.ReasoningContent != "" {
					m.stream, _ = updateStream(m.stream, StreamAppendMsg(choice.ReasoningContent), nil)
				}
				if choice.Content != "" {
					m.chat, _ = updateChat(m.chat, ChatAppendMsg(choice.Content), nil)
				}
			} else {
				m.timeline, _ = updateTimeline(m.timeline, TimelineAppendMsg("Warning: JSON parsed but 0 choices"), nil)
			}
		} else {
			m.timeline, _ = updateTimeline(m.timeline, TimelineAppendMsg("JSON Error (saved to debug_payload.json)"), nil)
		}
		m.timeline, _ = updateTimeline(m.timeline, TimelineAppendMsg(fmt.Sprintf("[%s] %s %s (%d)", msg.Proxy, msg.Method, msg.URL, msg.Status)), nil)
		cmds = append(cmds, readWS(m.wsConn))
	case WSErrMsg:
		m.timeline, _ = updateTimeline(m.timeline, TimelineAppendMsg("WS Error: " + msg.Error()), nil)
	case CommandSubmitMsg:
		// Echo command to chat and timeline
		cmdStr := string(msg)
		m.chat, _ = updateChat(m.chat, ChatAppendMsg("Operator: " + cmdStr), nil)
		m.timeline, _ = updateTimeline(m.timeline, TimelineAppendMsg("Command dispatched: " + cmdStr), nil)
		// TODO: Call m.client.SubmitCommand(cmdStr) async here
	}

	var newCmds []tea.Cmd
	m.stream, newCmds = updateStream(m.stream, msg, nil)
	if len(newCmds) > 0 { cmds = append(cmds, newCmds...) }
	
	m.timeline, newCmds = updateTimeline(m.timeline, msg, nil)
	if len(newCmds) > 0 { cmds = append(cmds, newCmds...) }
	
	m.chat, newCmds = updateChat(m.chat, msg, nil)
	if len(newCmds) > 0 { cmds = append(cmds, newCmds...) }
	
	m.command, newCmds = updateCommand(m.command, msg, nil, m.client)
	if len(newCmds) > 0 { cmds = append(cmds, newCmds...) }

	m.telemetry, newCmds = updateTelemetry(m.telemetry, msg, nil)
	if len(newCmds) > 0 { cmds = append(cmds, newCmds...) }

	return m, tea.Batch(cmds...)
}

func (m CockpitModel) View() string {
	header := m.styles.Header.Render(m.cfg.AppName)

	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.styles.StreamPane.Render(m.stream.View()),
		m.styles.TimelinePane.Render(m.timeline.View()),
	)

	middle := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.styles.ChatPane.Render(m.chat.View()),
		m.styles.TelemetryPane.Render(m.telemetry.View()),
	)

	bottom := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.styles.CommandPane.Render(m.command.View()),
		lipgloss.NewStyle().Render(""), // placeholder for future pane
	)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		top,
		middle,
		bottom,
	)

	return m.styles.Root.Render(body)
}

// helper to pass context later
func (m CockpitModel) context() context.Context {
	return context.Background()
}
