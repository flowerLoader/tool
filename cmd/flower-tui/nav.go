package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	zone "github.com/lrstanley/bubblezone"
)

type Item struct {
	Func func() string
	Name string
}

type navComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	SuppressKeyboard     bool
	SuppressMouseButtons bool
	SuppressMouseWheel   bool

	items    []Item
	onClick  func(int) tea.Cmd
	position int
	style    lipgloss.Style

	// render cache to reduce allocations
	elements []string
}

func NewNavComponent(
	style lipgloss.Style,
	items []Item,
	onClick func(int) tea.Cmd,
) *navComponent {
	return &navComponent{
		elements: make([]string, 2),
		items:    items,
		onClick:  onClick,
		position: 0,
		style:    lipgloss.NewStyle(),
	}
}

func (c *navComponent) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !c.SuppressKeyboard {
			switch msg.String() {
			case "up":
				c.position = max(c.position-1, 0)
			case "down":
				c.position = min(c.position+1, len(c.items)-1)
			}
		}

	case tea.MouseMsg:
		if !c.SuppressMouseWheel {
			switch msg.Button {
			case tea.MouseButtonWheelUp:
				c.position = max(c.position-1, 0)
			case tea.MouseButtonWheelDown:
				c.position = min(c.position+1, len(c.items)-1)
			}
		}

		if !c.SuppressMouseButtons {
			for i := 0; i <= len(c.items); i++ {
				if zone.Get(fmt.Sprintf("nav-%d", i)).InBounds(msg) {
					if msg.Action == tea.MouseActionPress || msg.Button == tea.MouseButtonLeft {
						c.position = i
						return c.onClick(c.position)
					}
				}
			}
		}
	}

	return nil
}

func (c *navComponent) Render(id, width, height int) string {
	if c.position == id {
		// UTF-8 right arrows: "▶", "➤", "➜", "➩", "➪", "➫", "➬", "➭", "➮", "➯", "➱", "➲", "➳", "➴", "➵", "➶", "➷", "➸", "➹", "➺", "➻", "➼", "➽", "➾"
		c.elements[0] = " ➽  "
	} else {
		c.elements[0] = "    "
	}

	if c.items[id].Func != nil {
		c.elements[1] = c.items[id].Func()
	} else {
		c.elements[1] = c.style.
			Height(height).
			Width(width - 4).
			Render(c.items[id].Name)
	}

	return zone.Mark(
		fmt.Sprintf("nav-%d", id),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			c.elements...,
		),
	)
}
