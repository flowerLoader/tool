package main

import "github.com/charmbracelet/lipgloss"

var theme = Theme{
	Default: Style{
		Foreground: "250",
	},
	Disabled: Style{
		Foreground: "243",
	},
	Border: Style{
		Foreground: "103",
	},
	Primary: Style{
		Foreground: "140",
	},
	Secondary: Style{
		Foreground: "97",
	},
	Error: Style{
		Foreground: "160",
	},
}

type Style struct {
	Background string `json:"background,omitempty"`
	Bold       bool   `json:"bold,omitempty"`
	Foreground string `json:"foreground,omitempty"`
}

type Theme struct {
	Default  Style `json:"default"`
	Disabled Style `json:"disabled"`

	// App
	Border Style `json:"border"`

	// Accent
	Primary   Style `json:"primary"`
	Secondary Style `json:"secondary"`

	// Status
	Error Style `json:"error"`
}

type StyleType string

const (
	DefaultStyle   StyleType = "default"
	DisabledStyle  StyleType = "disabled"
	BorderStyle    StyleType = "border"
	PrimaryStyle   StyleType = "primary"
	SecondaryStyle StyleType = "secondary"
	ErrorStyle     StyleType = "error"
)

func (t Theme) Gloss(name StyleType) lipgloss.Style {
	base := lipgloss.NewStyle().
		Bold(t.Default.Bold).
		Background(lipgloss.Color(t.Default.Background)).
		Foreground(lipgloss.Color(t.Default.Foreground))

	apply := func(s Style) {
		if s.Bold {
			base = base.Bold(true)
		}
		if s.Background != "" {
			base = base.Background(lipgloss.Color(s.Background))
		}
		if s.Foreground != "" {
			base = base.Foreground(lipgloss.Color(s.Foreground))
		}
	}

	switch name {
	case BorderStyle:
		if t.Border.Bold {
			base = base.Bold(true)
		}
		if t.Border.Background != "" {
			base = base.BorderBackground(lipgloss.Color(t.Border.Background)).
				MarginBackground(lipgloss.Color(t.Default.Background))
		}
		if t.Border.Foreground != "" {
			base = base.BorderForeground(lipgloss.Color(t.Border.Foreground))
		}
	case PrimaryStyle:
		apply(t.Primary)
	case SecondaryStyle:
		apply(t.Secondary)
	case DisabledStyle:
		apply(t.Disabled)
	case ErrorStyle:
		apply(t.Error)
	case DefaultStyle:
		// noop
	default:
		panic("unknown style: " + name)
	}

	return base
}
