package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/takacs/donkey/db"
)

type StatsModel struct {
	width, height int
	keys          keyMap
	help          help.Model
	name          string
}

func (m StatsModel) Init() tea.Cmd {
	return nil
}

func (m StatsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.MainMenu):
			path, err := db.GetDbPath("cards")
			if err != nil {
				fmt.Println("error getting db path")
			}
			return InitProject(path, m.width, m.height)
		default:
			fmt.Printf("default press quit %v \n", msg)
			return m, tea.Quit
		}
	}
	return m, tea.Batch(cmds...)
}

func (m StatsModel) View() string {
	helpView := m.help.View(m.keys)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		baseStyle.Render(m.name+"\n"+helpView))
}

func newStatsModel(width, height int) StatsModel {
	return StatsModel{
		width:  width,
		height: height,
		name:   "stats",
		help:   help.New(),
		keys:   keys,
	}
}
