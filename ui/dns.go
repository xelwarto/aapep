package ui

import (
	"aapep/util"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var renderWidth = 100

type DnsCmdModel struct {
	Res chan DnsCmdResMsg
	CMD struct {
		Clients  int
		Interval int
		Version  string
	}
	errors   int
	quitting bool
	average  time.Duration
}

type DnsCmdResMsg struct {
	Average time.Duration
	Errors  int
}

func waitForDNSCmdActivity(res chan DnsCmdResMsg) tea.Cmd {
	return func() tea.Msg {
		return DnsCmdResMsg(<-res)
	}
}

func (m DnsCmdModel) Init() tea.Cmd {
	return tea.Batch(waitForDNSCmdActivity(m.Res), tea.EnterAltScreen)
}

func (m DnsCmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case DnsCmdResMsg:
		m.average = msg.Average
		m.errors = msg.Errors
		return m, waitForDNSCmdActivity(m.Res)
	default:
		return m, nil
	}
}

func (m DnsCmdModel) View() string {
	w := lipgloss.Width
	bH1 := util.BannerHeader(m.CMD.Version, renderWidth)
	cH2 := util.CmdHeader.Render(fmt.Sprintf("Clients: %v", m.CMD.Clients))
	iH2 := util.CmdHeader.Render(fmt.Sprintf("Interval: %vms", m.CMD.Interval))
	tH2 := util.CmdHeader.Copy().Width(renderWidth - w(cH2) - w(iH2)).Render("DNS Test")
	bH2 := lipgloss.JoinHorizontal(lipgloss.Center, tH2, iH2, cH2)
	bA1 := fmt.Sprintf("\n Average RTT: %v", m.average)
	bE1 := fmt.Sprintf("\n Reported Errors: %v", util.ErrorMsgStyle.Render(fmt.Sprintf("%v", m.errors)))
	return bH1 + "\n" + bH2 + bA1 + bE1
}
