package command

import (
	"os"

	"github.com/renevo/mcutils/internal/command/clientcommand"
	"github.com/renevo/mcutils/internal/command/servercommand"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Execute(args []string) error {
	verboseLogging := false
	nocolorLogging := false
	jsonLogging := false

	rootCommand := &cobra.Command{
		Use:   "mcutils",
		Short: "Minecraft Utilities",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(os.Stdout)

			if jsonLogging {
				logrus.SetFormatter(&logrus.JSONFormatter{})
			} else {
				logrus.SetFormatter(&logrus.TextFormatter{
					DisableColors: nocolorLogging,
				})
			}

			if verboseLogging {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				logrus.SetLevel(logrus.InfoLevel)
			}

			logrus.WithField("command", cmd.Use).Debug("Command PersistentPreRunE")

			return nil
		},
	}

	rootCommand.PersistentFlags().BoolVarP(&verboseLogging, "verbose", "v", false, "verbose output")
	rootCommand.PersistentFlags().BoolVarP(&jsonLogging, "json", "j", false, "output logging as json")
	rootCommand.PersistentFlags().BoolVar(&nocolorLogging, "no-color", false, "disable colorized output")

	// add commands here:
	rootCommand.AddCommand(
		clientcommand.New(),
		servercommand.New(),
	)

	// execute
	rootCommand.SetArgs(args)
	return rootCommand.Execute()
}
