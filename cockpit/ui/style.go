package ui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/pqr-info/pqr-info-swarm/cockpit/internal/config"
)

type Styles struct {
	Root        lipgloss.Style
	Header      lipgloss.Style
	PaneBorder  lipgloss.Style
	StreamPane  lipgloss.Style
	TimelinePane lipgloss.Style
	ChatPane    lipgloss.Style
	CommandPane lipgloss.Style
	TelemetryPane lipgloss.Style
}

func NewStyles(cfg config.TenantConfig) Styles {
	bg := cfg.Theme.Background
	primary := cfg.Theme.Primary
	secondary := cfg.Theme.Secondary

	root := lipgloss.NewStyle().
		Background(lipgloss.Color(bg))

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color(primary)).
		Bold(true)

	border := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(secondary)).
		Padding(0, 1)

	return Styles{
		Root:        root,
		Header:      header,
		PaneBorder:  border,
		StreamPane:  border,
		TimelinePane: border,
		ChatPane:    border,
		CommandPane: border,
		TelemetryPane: border,
	}
}
