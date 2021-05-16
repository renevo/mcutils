package clientcommand

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	clientCommand := &cobra.Command{
		Use:   "client",
		Short: "Runs a minecraft client",
		Long:  `Will download, configure, and run a Minecraft client. Without any configuration, this will be a default vanilla minecraft client`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.WithFields(logrus.Fields{"client": "1.16"}).Info("Initialize the client here....")

			return nil
		},
	}

	return clientCommand
}
