package util

import "github.com/charmbracelet/lipgloss"

var (
	ErrorMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FC0000"))

	bannerHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#679436"))

	bannerHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ebf2fa")).
				Background(lipgloss.Color("#679436")).
				Padding(0, 1)

	versionHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ebf2fa")).
				Background(lipgloss.Color("#05668d")).
				Padding(0, 1)

	CmdHeader = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ebf2fa")).
			Background(lipgloss.Color("#05668d")).Padding(0, 1)
)
