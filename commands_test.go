package main

import (
	"testing"

	"github.com/urfave/cli"
)

func TestCommands(t *testing.T) {
	var cs, subcs []cli.Command

	for _, c := range Commands {
		if len(c.Subcommands) == 0 {
			cs = append(cs, c)
		} else {
			for _, sc := range c.Subcommands {
				cs = append(cs, sc)
			}
			subcs = append(subcs, c)
		}
	}

	// Flags が存在する場合は ArgsUsage を必須とする
	for _, c := range cs {
		if len(c.Flags) > 0 && c.ArgsUsage == "" {
			t.Errorf("%s: cli.Command.ArgsUsage should not be empty.", c.Name)
		}
	}

	for _, sc := range subcs {
		if sc.Action == nil {
			if sc.Usage == "" {
				t.Errorf("%s: Neither .Description nor .Usage should be empty", sc.Name)

			}
		}
	}
}
