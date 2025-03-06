package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/zachdehooge/Weather-Dashboard/cmd"
)

func main() {

	var initial string

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Margin(1, 0).
		Bold(true)

	title := titleStyle.Render("Weather Dashboard")

	fmt.Println(title)

	initialText := "\n(1)All Alerts\n(2)State Alerts\n(3)day-forecast\n(4)exit\n\nChoose An Option: "

	// Create a Lipgloss style with a border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		Margin(0).
		BorderForeground(lipgloss.Color("#FF69B4")) // Optional: Color the border

	// Apply the style to the text
	styledText := borderStyle.Render(initialText)

	// Print the styled text+
	fmt.Println(styledText)

	// TODO: Add Weather Trends Analyzer and Convective Outlook Options

	fmt.Scanln(&initial)

	if initial == "1" {
		cmd.AllAlerts()
		main()
	} else if initial == "2" {
		cmd.StateAlerts()
		main()
	} else if initial == "3" {
		cmd.Forecast()
		main()
	} else {
		fmt.Println("\nExiting Application")
	}
}
