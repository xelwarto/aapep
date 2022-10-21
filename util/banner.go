package util

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var helpLogo = `
   _                          
  /_\   __ _ _ __   ___ _ __  
 //_\\ / _' | '_ \ / _ \ '_ \ 
/  _  \ (_| | |_) |  __/ |_) |
\_/ \_/\__,_| .__/ \___| .__/ 
            |_|        |_| v%v   
`
var headerLogo = `//\apep - Stress Tester`

func BannerHelp(v string) string {
	logo := fmt.Sprintf(helpLogo, v)
	return bannerHelpStyle.Render(logo)
}

func BannerHeader(v string, size int) string {
	w := lipgloss.Width
	ver := versionHeaderStyle.Render("v" + v)
	return lipgloss.JoinHorizontal(lipgloss.Top, bannerHeaderStyle.Copy().Width(size-w(ver)).Render(headerLogo), ver)
}
