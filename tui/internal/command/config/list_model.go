package config

import (
	"sync"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/fressive/pocman/tui/internal/conf"
)

type listModel struct {
	styles        styles
	darkBG        bool
	width, height int
	once          *sync.Once
	list          list.Model
	delegateKeys  *delegateKeyMap
}

func (m *configModel) updateListProperties() {
	// Update list size.
	h, v := m.styles.list.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)

	// Update the model and list styles.
	m.styles = newStyles()
	m.list.Styles.Title = m.styles.title
}

func (m configModel) updateList(msg tea.Msg) (configModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.darkBG = msg.IsDark()
		m.updateListProperties()
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil

	case itemMsg:
		m.selected = msg.Item
		m.textInput.SetValue(*m.selected.value)
		return m, nil
	}

	switch msg.(type) {
	case tea.KeyPressMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m configModel) listView() tea.View {
	return tea.NewView(m.styles.list.Render(m.list.View()))
}

func (m *configModel) initList() {
	delegateKeys := newDelegateKeyMap()

	items := []list.Item{
		item{title: "Endpoint", value: &conf.TUIConfig.Server.Endpoint},
		item{title: "Token", value: &conf.TUIConfig.Server.Token},
	}

	// Setup list.
	delegate := m.newItemDelegate(delegateKeys)
	confList := list.New(items, delegate, 0, 0)
	confList.SetShowStatusBar(false)
	confList.Title = "Configuration"
	confList.Styles.Title = m.styles.title

	m.list = confList
	m.delegateKeys = delegateKeys
}
