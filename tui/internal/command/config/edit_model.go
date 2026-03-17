package config

import (
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/fressive/pocman/tui/internal/conf"
)

type editModel struct {
	selected  *item
	textInput textinput.Model
}

type editKeyMap struct {
	Confirm key.Binding
	Reset   key.Binding
	Quit    key.Binding
}

func (k editKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Reset, k.Quit}
}

func (k editKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Confirm, k.Reset, k.Quit},
	}
}

var editKeys = editKeyMap{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit w/o saving"),
	),
}

func (m configModel) updateEdit(msg tea.Msg) (configModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.saveConfig()
			m.selected = nil
			return m, nil
		case "r":
			m.textInput.SetValue(*m.selected.value)
			return m, nil
		case "q", "ctrl+c", "esc":
			m.selected = nil
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m configModel) editView() tea.View {
	var c *tea.Cursor

	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
	}

	help := help.New()
	helpView := help.View(editKeys)

	str := m.styles.edit.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			m.styles.title.Render("Edit"),
			"\n",
			m.styles.editTitle.Render(fmt.Sprintf("%s", m.selected.title)),
			m.textInput.View(),
			"\n",
			helpView,
		),
	)

	v := tea.NewView(str)
	v.Cursor = c
	return v
}

func (m *configModel) saveConfig() {
	*m.selected.value = m.textInput.Value()

	path, err := conf.DefaultFilePath()

	if err != nil {
		panic(err)
	}

	err = conf.TUIConfig.Save(path)

	if err != nil {
		panic(err)
	}
}

func (m *configModel) initEdit() {
	ti := textinput.New()
	ti.Focus()
	ti.SetWidth(200)
	m.textInput = ti
}
