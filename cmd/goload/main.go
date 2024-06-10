package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"goload/internal/configs"
	"goload/internal/wiring"
	"log"
)

var (
	version    string
	commitHash string
)

const (
	flagConfigFilePath = "config-file-path"
)

func server() *cobra.Command {
	command := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFilePath, err := cmd.Flags().GetString(flagConfigFilePath)

			if err != nil {
				return err
			}

			app, cleanup, err := wiring.InitializeGRPCServer(configs.ConfigFilePath(configFilePath))
			if err != nil {
				return err
			}

			defer cleanup()

			return app.Start()
		},
	}

	command.Flags().String(flagConfigFilePath, "", "If provided, will use the provide config file")
	return command
}

func main() {
	rootCommand := &cobra.Command{
		Version: fmt.Sprintf("%s-%s", version, commitHash),
	}
	rootCommand.AddCommand(
		server(),
	)

	if err := rootCommand.Execute(); err != nil {
		log.Panicln(err)
	}
}
