package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#04B575"))

func Before(c *cli.Context) error {
	profile := c.String("profile")
	region := c.String("region")

	fmt.Println(style.Render("Profile: "+profile, "Region: "+region))

	return nil
}
