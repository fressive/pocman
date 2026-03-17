package config

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	list      lipgloss.Style
	edit      lipgloss.Style
	title     lipgloss.Style
	editTitle lipgloss.Style
}

func newStyles() styles {
	return styles{
		list: lipgloss.NewStyle().
			Padding(1, 1),
		edit: lipgloss.NewStyle().
			Padding(1, 3),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		editTitle: lipgloss.NewStyle().
			Bold(true),
	}
}

type item struct {
	title string
	value *string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return *i.value }
func (i item) FilterValue() string { return i.title + *i.value }

type configModel struct {
	listModel
	editModel
}

func (m configModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.selected == nil {
		// render the list view
		m, cmd = m.updateList(msg)
	} else {
		// render the edit view
		m, cmd = m.updateEdit(msg)
	}

	return m, cmd
}

func (m configModel) View() tea.View {
	var v tea.View

	if m.selected != nil {
		v = m.editView()
	} else {
		v = m.listView()
	}

	v.AltScreen = true
	return v
}

func initialModel() configModel {
	m := configModel{}
	m.styles = newStyles()

	m.initList()
	m.initEdit()

	return m
}
