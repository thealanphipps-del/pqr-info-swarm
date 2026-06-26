package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thealanphipps-del/pqr"
)

// Messages for tea.Model loop
type tickMsg time.Time
type refetchMsg struct {
	Timeline []string
	Logs     []string
}

type Model struct {
	textInput    textinput.Model
	Logs         []string
	Timeline     []string
	ChatHistory  []string
	ActiveAgent  string
	Paused       bool
	DashboardURL string
	APIURL       string
	Client       *pqr.Client
	StatusMsg    string
	Quitting     bool
}

func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Type command (/chat <id>, /ask <msg>, /take <id>, /assign <id> <agent>, /swarm kill, /exit)..."
	ti.Focus()
	ti.CharLimit = 250
	ti.Width = 80

	dashboardURL := os.Getenv("SWEN_DASHBOARD_URL")
	if dashboardURL == "" {
		dashboardURL = "http://127.0.0.1:7777"
	}
	apiURL := os.Getenv("PQR_API_URL")
	if apiURL == "" {
		apiURL = "http://127.0.0.1:8196"
	}

	return Model{
		textInput:    ti,
		Logs:         []string{"Swarm stream initialized.", "Observing node on " + dashboardURL},
		Timeline:     []string{"No timeline events loaded."},
		ChatHistory:  []string{"No active chat session. Use /chat <agent-id> to begin."},
		ActiveAgent:  "none",
		Paused:       false,
		DashboardURL: dashboardURL,
		APIURL:       apiURL,
		Client:       pqr.NewClient(apiURL),
		StatusMsg:    "Ready. Press Enter to submit commands or refresh.",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.refetchCmd(),
		m.tickCmd(),
	)
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) refetchCmd() tea.Cmd {
	return func() tea.Msg {
		timeline := []string{}
		logs := []string{}

		// 1. Fetch Timeline
		resp, err := http.Get(m.DashboardURL + "/timeline")
		if err == nil {
			defer resp.Body.Close()
			var events []map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&events); err == nil {
				start := len(events) - 6
				if start < 0 {
					start = 0
				}
				for i := len(events) - 1; i >= start; i-- {
					evt := events[i]
					tStr, _ := evt["timestamp"].(string)
					t, _ := time.Parse(time.RFC3339, tStr)
					evtType, _ := evt["type"].(string)
					data, _ := evt["data"].(map[string]interface{})

					summary := ""
					if evtType == "snapshot" {
						summary = fmt.Sprintf("[%s] SNAPSHOT: Genesis=%v", t.Format("15:04:05"), data["is_genesis"])
					} else if evtType == "journal" {
						summary = fmt.Sprintf("[%s] JOURNAL: %s", t.Format("15:04:05"), data["action"])
					} else {
						summary = fmt.Sprintf("[%s] ERROR: %s", t.Format("15:04:05"), data["signature"])
					}
					timeline = append(timeline, summary)
				}
			}
		}

		// 2. Fetch Journal for logs
		resp2, err2 := http.Get(m.DashboardURL + "/journal")
		if err2 == nil {
			defer resp2.Body.Close()
			var journal []map[string]interface{}
			if err := json.NewDecoder(resp2.Body).Decode(&journal); err == nil {
				for i := len(journal) - 1; i >= 0; i-- {
					j := journal[i]
					action, _ := j["action"].(string)
					tStr, _ := j["timestamp"].(string)
					t, _ := time.Parse(time.RFC3339, tStr)
					logs = append(logs, fmt.Sprintf("[%s] ACTION: %s", t.Format("15:04:05"), action))
				}
			}
		}

		return refetchMsg{Timeline: timeline, Logs: logs}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			input := strings.TrimSpace(m.textInput.Value())
			m.textInput.SetValue("")
			if input == "" {
				return m, m.refetchCmd()
			}

			// Handle subcommands
			if strings.HasPrefix(input, "/") {
				parts := strings.SplitN(input, " ", 3)
				c := parts[0]

				switch c {
				case "/exit", "/quit":
					m.Quitting = true
					return m, tea.Quit
				case "/pause":
					m.Paused = true
					m.Logs = append(m.Logs, "Ingestion paused by operator.")
				case "/resume":
					m.Paused = false
					m.Logs = append(m.Logs, "Ingestion resumed by operator.")
				case "/swarm":
					if len(parts) > 1 && (parts[1] == "kill" || parts[1] == "stop") {
						m.Logs = append(m.Logs, "Triggering swarm kill-switch...")
						ExecKillSwitch()
						m.Logs = append(m.Logs, "Legacy processes killed successfully.")
					} else {
						m.Logs = append(m.Logs, "Usage: /swarm kill")
					}
				case "/chat":
					if len(parts) < 2 {
						m.Logs = append(m.Logs, "Usage: /chat <agent-id>")
						break
					}
					m.ActiveAgent = parts[1]
					m.Logs = append(m.Logs, "Switched chat focus to Agent: "+parts[1])
					m.LoadConversation()
				case "/ask":
					if len(parts) < 2 {
						m.Logs = append(m.Logs, "Usage: /ask <message>")
						break
					}
					msgVal := parts[1]
					if len(parts) > 2 {
						msgVal = parts[1] + " " + parts[2]
					}
					if m.ActiveAgent == "none" {
						m.Logs = append(m.Logs, "Select an agent first using /chat <agent-id>")
						break
					}
					m.SendChatMessage(msgVal)
				case "/take":
					if len(parts) < 2 {
						m.Logs = append(m.Logs, "Usage: /take <ticket-id>")
						break
					}
					m.Logs = append(m.Logs, "Taking ticket: "+parts[1])
					err := m.Client.UpdateTicketExtended(context.Background(), parts[1], "", "", "operator", "operator", "", "")
					if err != nil {
						m.Logs = append(m.Logs, "Error: "+err.Error())
					} else {
						m.Logs = append(m.Logs, "Ticket successfully taken.")
					}
				case "/assign":
					if len(parts) < 3 {
						m.Logs = append(m.Logs, "Usage: /assign <ticket-id> <agent-id>")
						break
					}
					m.Logs = append(m.Logs, fmt.Sprintf("Assigning ticket %s to %s...", parts[1], parts[2]))
					err := m.Client.UpdateTicketExtended(context.Background(), parts[1], "", "", "operator", parts[2], "", "")
					if err != nil {
						m.Logs = append(m.Logs, "Error: "+err.Error())
					} else {
						m.Logs = append(m.Logs, "Ticket successfully assigned.")
					}
				default:
					m.Logs = append(m.Logs, "Unknown command: "+c)
				}
			} else {
				if m.ActiveAgent != "none" {
					m.SendChatMessage(input)
				} else {
					m.Logs = append(m.Logs, "No active agent. Use /chat <agent-id> to talk, or prefix with '/' for commands.")
				}
			}
			return m, m.refetchCmd()
		}

	case tickMsg:
		if !m.Paused {
			return m, tea.Batch(m.refetchCmd(), m.tickCmd())
		}
		return m, m.tickCmd()

	case refetchMsg:
		if len(msg.Timeline) > 0 {
			m.Timeline = msg.Timeline
		}
		// Merges unique new log strings
		for _, l := range msg.Logs {
			exists := false
			for _, current := range m.Logs {
				if current == l {
					exists = true
					break
				}
			}
			if !exists {
				m.Logs = append(m.Logs, l)
			}
		}
		if len(m.Logs) > 30 {
			m.Logs = m.Logs[len(m.Logs)-30:]
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *Model) LoadConversation() {
	resp, err := http.Get(fmt.Sprintf("%s/REST/2.0/agent/%s/conversation", m.APIURL, m.ActiveAgent))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var logs []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return
	}

	m.ChatHistory = []string{}
	for _, l := range logs {
		content, _ := l["content"].(string)
		m.ChatHistory = append(m.ChatHistory, content)
	}

	if len(m.ChatHistory) == 0 {
		m.ChatHistory = []string{"No conversation history with Agent " + m.ActiveAgent}
	}
}

func (m *Model) SendChatMessage(msg string) {
	m.Logs = append(m.Logs, fmt.Sprintf("Operator -> Agent %s: %s", m.ActiveAgent, msg))

	payload := map[string]interface{}{
		"sender":  "operator",
		"message": msg,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(fmt.Sprintf("%s/REST/2.0/agent/%s/message", m.APIURL, m.ActiveAgent), "application/json", strings.NewReader(string(body)))
	if err != nil {
		m.Logs = append(m.Logs, "Error sending message: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		agentResp, _ := result["response"].(string)
		m.Logs = append(m.Logs, agentResp)
	}
	m.LoadConversation()
}

func (m Model) View() string {
	if m.Quitting {
		return "Shutdown signal received. Exiting Cockpit.\n"
	}

	// Layout and Styles using lipgloss
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#00f0ff")).
		Padding(0, 2).
		MarginBottom(1)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4a5568")).
		Padding(0, 1)

	header := headerStyle.Render("=== SWEN SWARM OPERATOR COCKPIT & CONTROL CENTER ===")

	// Top stream box
	streamContent := ""
	logStart := len(m.Logs) - 6
	if logStart < 0 {
		logStart = 0
	}
	for i := logStart; i < len(m.Logs); i++ {
		streamContent += fmt.Sprintf("• %s\n", m.Logs[i])
	}
	streamView := borderStyle.
		BorderForeground(lipgloss.Color("#10b981")).
		Width(82).
		Render(fmt.Sprintf("Live Swarm Stream:\n%s", streamContent))

	// Left and Right Columns
	leftContent := ""
	for i, entry := range m.Timeline {
		if i < 6 {
			leftContent += fmt.Sprintf("%s\n", entry)
		}
	}
	timelineView := borderStyle.
		BorderForeground(lipgloss.Color("#3182ce")).
		Width(39).
		Height(8).
		Render(fmt.Sprintf("Timeline Snapshots:\n%s", leftContent))

	rightContent := ""
	for i, entry := range m.ChatHistory {
		if i < 6 {
			rightContent += fmt.Sprintf("%s\n", entry)
		}
	}
	chatView := borderStyle.
		BorderForeground(lipgloss.Color("#d53f8c")).
		Width(39).
		Height(8).
		Render(fmt.Sprintf("Chat (Agent: %s):\n%s", m.ActiveAgent, rightContent))

	columns := lipgloss.JoinHorizontal(lipgloss.Top, timelineView, chatView)

	footer := fmt.Sprintf("\n%s\n%s", m.StatusMsg, m.textInput.View())

	return fmt.Sprintf("%s\n%s\n%s\n%s\n", header, streamView, columns, footer)
}

func StartTUI() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("TUI startup error: %v", err)
		os.Exit(1)
	}
}

func ExecKillSwitch() {
	if runtime.GOOS == "windows" {
		exec.Command("taskkill", "/F", "/IM", "python.exe").Run()
		exec.Command("taskkill", "/F", "/IM", "gemma.exe").Run()
		exec.Command("taskkill", "/F", "/IM", "inference.exe").Run()
	} else {
		exec.Command("pkill", "-f", "gemma").Run()
		exec.Command("pkill", "-f", "inference").Run()
		exec.Command("pkill", "-f", "swarm").Run()
		exec.Command("sudo", "systemctl", "stop", "swarm-agent.service").Run()
		exec.Command("sudo", "systemctl", "stop", "swarm-runner.service").Run()
		exec.Command("sudo", "systemctl", "stop", "inference-worker.service").Run()
	}
}
