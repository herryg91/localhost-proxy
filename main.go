package main

import (
	"os"

	"github.com/herryg91/localhost-proxy/handler/cli"
	cli_config "github.com/herryg91/localhost-proxy/handler/cli/config"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "localhost-proxy", Short: "lprx", Long: "localhost-proxy"}
	rootCmd.AddCommand(cli.NewCmdStatus().Command)
	rootCmd.AddCommand(cli_config.New().Command)
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println(err)
		os.Exit(1)
	}

	return
}
