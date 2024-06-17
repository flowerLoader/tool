package main

import "github.com/charmbracelet/lipgloss"

var theme = &Theme{
	Name: "Flower Dark",
	Styles: map[StyleType]Style{
		DefaultStyle: {
			Foreground: "250",
		},
		DisabledStyle: {
			Foreground: "243",
		},
		BorderStyle: {
			Foreground: "103",
		},
		PrimaryStyle: {
			Foreground: "140",
		},
		SecondaryStyle: {
			Foreground: "97",
		},
		ErrorStyle: {
			Foreground: "160",
		},
	},
}

type Style struct {
	Background string `json:"background,omitempty"`
	Bold       bool   `json:"bold,omitempty"`
	Foreground string `json:"foreground,omitempty"`
}

type Theme struct {
	Name   string              `json:"name,omitempty"`
	Styles map[StyleType]Style `json:"styles,omitempty"`
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
		Bold(t.Styles[DefaultStyle].Bold).
		Background(lipgloss.Color(t.Styles[DefaultStyle].Background)).
		Foreground(lipgloss.Color(t.Styles[DefaultStyle].Foreground))

	s := t.Styles[name]
	if s.Bold {
		base = base.Bold(true)
	}

	switch name {
	case BorderStyle:
		if s.Background != "" {
			base = base.BorderBackground(lipgloss.Color(s.Background)).
				MarginBackground(lipgloss.Color(s.Background))
		}
		if s.Foreground != "" {
			base = base.BorderForeground(lipgloss.Color(s.Foreground))
		}
	case DefaultStyle:
		// noop
	default:
		if s.Background != "" {
			base = base.Background(lipgloss.Color(s.Background))
		}
		if s.Foreground != "" {
			base = base.Foreground(lipgloss.Color(s.Foreground))
		}
	}

	return base
}
