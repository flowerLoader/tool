package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimaryMain = lipgloss.Color("176")
	ColorPrimaryDark = lipgloss.Color("96")
	ColorDisabled    = lipgloss.Color("243")

	TextDark     = lipgloss.NewStyle().Foreground(ColorPrimaryDark)
	TextDisabled = lipgloss.NewStyle().Foreground(ColorDisabled)
	TextMain     = lipgloss.NewStyle().Foreground(ColorPrimaryMain)
)
