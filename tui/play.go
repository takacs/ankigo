package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type PlayModel struct {
	width, height int
	keys          keyMap
	help          help.Model
	name          string
}

func (m PlayModel) Init() tea.Cmd {
	return nil
}

func (m PlayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.MainMenu):
			return InitProject(m.width, m.height)
		default:
			fmt.Printf("default press quit %v \n", msg)
			return m, tea.Quit
		}
	}
	return m, tea.Batch(cmds...)
}
func (m PlayModel) View() string {
	helpView := m.help.View(m.keys)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		baseStyle.Render(m.name+"\n"+helpView))
}

func newPlayModel(width, height int) PlayModel {
	return PlayModel{
		width:  width,
		height: height,
		name:   "play",
		help:   help.New(),
		keys:   keys,
	}
}
