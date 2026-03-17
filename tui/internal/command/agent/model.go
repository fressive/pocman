package agent

import (
	"fmt"
	"math"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/fressive/pocman/common/pkg/model"
	"github.com/samber/lo"
)

type agentModel struct {
	table table.Model
	err   error
}

type agentMsg []model.Agent

type errMsg struct{ err error }

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func initAgentModel() agentModel {
	return agentModel{}
}

func (m agentModel) Init() tea.Cmd {
	return nil
}

func (m agentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case agentMsg:
		columns := []table.Column{
			{Title: "ID", Width: 10},
			{Title: "Status", Width: 10},
			{Title: "Uptime", Width: 10},
			{Title: "CPU", Width: 10},
			{Title: "RAM", Width: 20},
			{Title: "Load", Width: 20},
			{Title: "Tasks", Width: 10},
		}

		rows := lo.Map(msg, func(a model.Agent, _ int) table.Row {
			var online string
			var uptime string

			if a.Online {
				online = "Online"
				uptimeDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", math.Round(a.Uptime)))
				uptime = uptimeDuration.String()
			} else {
				online = "Offline"
				uptime = "N/A"
			}

			return table.Row{
				a.AgentID,
				online,
				uptime,
				fmt.Sprintf("%.0f%%", a.CPUUsage),
				fmt.Sprintf("%dM/%dM", (a.RAMTotal-a.RAMAvailable)/1024/1024, a.RAMTotal/1024/1024),
				fmt.Sprintf("%.2f %.2f %.2f", a.Load1, a.Load5, a.Load15),
			}

		})

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(7),
			table.WithWidth(100),
		)

		m.table, cmd = t.Update(msg)
	case errMsg:
		m.err = msg.err
		cmd = tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			cmd = tea.Quit
		}
	}

	return m, cmd
}

func (m agentModel) View() tea.View {
	if m.err != nil {
		return tea.NewView(fmt.Sprintf("\nError: %v\n\n", m.err))
	}

	return tea.NewView(baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n")
}
