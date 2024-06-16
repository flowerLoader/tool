package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ANSIBackground = "0"
	ANSIBorder     = "103"
	ANSIPrimary    = "140"
	ANSISecondary  = "97"
	ANSIDisabled   = "243"
	ANSIError      = "160"

	ColorBorder    = lipgloss.NewStyle().BorderForeground(lipgloss.Color(ANSIBorder))
	ColorPrimary   = lipgloss.NewStyle().Foreground(lipgloss.Color(ANSIPrimary))
	ColorSecondary = lipgloss.NewStyle().Foreground(lipgloss.Color(ANSISecondary))
	ColorDisabled  = lipgloss.NewStyle().Foreground(lipgloss.Color(ANSIDisabled))
	ColorError     = lipgloss.NewStyle().Foreground(lipgloss.Color(ANSIError))
)
