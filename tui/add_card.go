package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	ddb "github.com/takacs/donkey/db"
)

const (
	front = iota
	back
	deck
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
	green    = lipgloss.Color("#55FF33")
)

const (
	formWidth = 50
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
	insertedStyle = lipgloss.NewStyle().Foreground(green)
)

type AddCardModel struct {
	width, height int
	inputs        []textinput.Model
	focus         int
	inserted      string
	keys          keyMap
	help          help.Model
	name          string
}

func (m AddCardModel) Init() tea.Cmd {
	return nil
}

func (m AddCardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Exit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Tab):
			m.nextFocus()
		case key.Matches(msg, m.keys.Enter):
			m.submitCard()
		}
	}
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m AddCardModel) View() string {
	helpView := m.help.View(m.keys)

	formView := fmt.Sprintf(
		` Add Card

 %s

 %s


 %s

 %s


 %s

 %s

 %s

 %s
`,
		inputStyle.Width(formWidth).Render("Front"),
		m.inputs[front].View(),
		inputStyle.Width(formWidth).Render("Back"),
		m.inputs[back].View(),
		inputStyle.Width(formWidth).Render("Deck"),
		m.inputs[deck].View(),
		insertedStyle.Render(m.inserted),
		continueStyle.Render("Continue ->"),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		baseStyle.Render(formView+"\n"+helpView))
}

func (m *AddCardModel) nextFocus() {
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + 1) % 3
	m.inputs[m.focus].Focus()
}

func (m *AddCardModel) submitCard() {
	c, err := ddb.OpenDb(ddb.SetupPath())
	if err != nil {
		fmt.Println("couldn't open db")
	}
	defer c.Db.Close()
	err = c.Insert(m.inputs[front].Value(), m.inputs[back].Value(), m.inputs[deck].Value())
	if err != nil {
		fmt.Println("couldn't instert card")
	}

	m.inputs = defaultInputs(0)
	m.inserted = "Inserted!"
}

func defaultInputs(focus int) []textinput.Model {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[front] = textinput.New()
	inputs[front].Placeholder = "donkey"
	inputs[front].Width = formWidth
	inputs[front].Prompt = "> "

	inputs[back] = textinput.New()
	inputs[back].Placeholder = "best app ever"
	inputs[back].CharLimit = 1000
	inputs[back].Width = formWidth
	inputs[back].Prompt = "> "

	inputs[deck] = textinput.New()
	inputs[deck].Placeholder = "default"
	inputs[deck].CharLimit = 100
	inputs[deck].Width = formWidth
	inputs[deck].Prompt = "> "

	inputs[focus].Focus()

	return inputs

}

func newAddCardModel(width, height int) AddCardModel {
	focus := 0
	inputs := defaultInputs(focus)

	return AddCardModel{
		width:    width,
		height:   height,
		inputs:   inputs,
		focus:    focus,
		inserted: "",
		name:     "add card",
		help:     help.New(),
		keys:     keys,
	}
}
