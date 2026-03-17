package test

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/fressive/pocman/tui/internal/api"
	"github.com/fressive/pocman/tui/internal/conf"
)

type testModel struct {
	success bool
	err     error
}

type successMsg struct{}
type errMsg struct{ error }

func checkConnection() tea.Msg {
	ctx := context.Background()

	c, err := api.GetClient()

	if err != nil {
		return errMsg{err}
	}

	err = c.Ping(ctx)

	if err != nil {
		return errMsg{err}
	}

	return successMsg{}
}

func (m testModel) Init() tea.Cmd {
	return checkConnection
}

func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.error
		return m, tea.Quit
	case successMsg:
		m.success = true
		return m, tea.Quit
	}

	return m, nil
}

func (m testModel) View() tea.View {
	s := fmt.Sprintf(fmt.Sprintf("Checking connection to server %s...\n\n", conf.TUIConfig.Server.Endpoint))

	if m.err != nil {
		s += fmt.Sprintf("❌ Error occured, check your configuration.\n\n%s", m.err)
	} else if m.success {
		s += "✅ Connect successfully."
	}

	return tea.NewView(s + "\n")
}

func initTestModel() testModel {
	return testModel{}
}
