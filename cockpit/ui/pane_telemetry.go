package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type TelemetryData struct {
	FreeDiskGB  float64
	LastCleanup time.Time
	Events      []string
	
	// Hardware Metrics
	CPUHistory  []float64
	RAMHistory  []float64
	GPU0History []float64 // NVIDIA
	GPU1History []float64 // Intel
	GPU2History []float64 // Yoga NPU
}

type TelemetryPaneModel struct {
	viewport viewport.Model
	data     TelemetryData
}

func NewTelemetryPaneModel() TelemetryPaneModel {
	v := viewport.New(40, 10)
	v.SetContent("Telemetry\n\nWaiting for data...")

	return TelemetryPaneModel{
		viewport: v,
		data: TelemetryData{
			FreeDiskGB:  0,
			LastCleanup: time.Time{},
			Events:      []string{},
			CPUHistory:  []float64{},
			RAMHistory:  []float64{},
			GPU0History: []float64{},
			GPU1History: []float64{},
			GPU2History: []float64{},
		},
	}
}

func (m TelemetryPaneModel) Update(msg tea.Msg) (TelemetryPaneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case TelemetryUpdateMsg:
		m.data = msg.Data
		m.viewport.SetContent(renderTelemetry(m.data))
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m TelemetryPaneModel) View() string {
	return m.viewport.View()
}

type TelemetryUpdateMsg struct {
	Data TelemetryData
}

func renderTelemetry(d TelemetryData) string {
	out := "Telemetry\n\n"

	out += fmt.Sprintf("CPU Load:      %s\n", renderSparkline(d.CPUHistory))
	out += fmt.Sprintf("RAM Usage:     %s\n", renderSparkline(d.RAMHistory))
	out += fmt.Sprintf("GPU0 (NVIDIA): %s\n", renderSparkline(d.GPU0History))
	out += fmt.Sprintf("GPU1 (Intel):  %s\n", renderSparkline(d.GPU1History))
	out += fmt.Sprintf("GPU2 (NPU):    %s\n\n", renderSparkline(d.GPU2History))

	out += fmt.Sprintf("Free Disk: %.2f GB\n", d.FreeDiskGB)

	if !d.LastCleanup.IsZero() {
		out += fmt.Sprintf("Last Clean: %s\n", d.LastCleanup.Format(time.RFC822))
	} else {
		out += "Last Clean: (none)\n"
	}

	out += "\nRecent Events:\n"
	if len(d.Events) == 0 {
		out += "  (no events)\n"
	} else {
		for _, e := range d.Events {
			out += "  • " + e + "\n"
		}
	}

	return out
}

func renderSparkline(data []float64) string {
	if len(data) == 0 {
		return "(no data)"
	}
	// " ▂▃▄▅▆▇█"
	chars := []rune(" ▂▃▄▅▆▇█")
	out := ""
	
	// only show last 25 items for sparkline
	start := 0
	if len(data) > 25 {
		start = len(data) - 25
	}

	for _, val := range data[start:] {
		idx := int((val / 100.0) * float64(len(chars)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(chars) {
			idx = len(chars) - 1
		}
		out += string(chars[idx])
	}
	
	// append the latest numerical value
	out += fmt.Sprintf(" [%.1f%%]", data[len(data)-1])
	return out
}

func updateTelemetry(m TelemetryPaneModel, msg tea.Msg, cmds []tea.Cmd) (TelemetryPaneModel, []tea.Cmd) {
	var cmd tea.Cmd
	m, cmd = m.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, cmds
}

