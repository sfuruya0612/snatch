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
	p := c.String("profile")
	r := c.String("region")

	fmt.Println(style.Render("Profile: "+p, "Region: "+r))

	return nil
}
