package cli_config

import (
	"github.com/spf13/cobra"
)

type CmdConfig struct {
	*cobra.Command
}

func New() *CmdConfig {
	c := &CmdConfig{}
	c.Command = &cobra.Command{
		Use:   "config",
		Short: "Configure server & preferences",
		Long:  "Configure server & preferences",
	}

	c.AddCommand(newConfigEdit().Command)
	c.AddCommand(newConfigGet().Command)
	return c
}
