package minecraft

import (
	"context"

	"github.com/pkg/errors"
	"github.com/portcullis/application"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/cnc"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/mcserver"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/pubsub"
	"github.com/renevo/mcutils/pkg/minecraft"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func serverCommands() []*cobra.Command {
	configFile := "./minecraft.hcl"
	srv := minecraft.Default()

	commands := []*cobra.Command{}

	commands = append(commands, &cobra.Command{
		Use:   "validate",
		Short: "Validates configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := application.New("mcutils", "1.0.0",
				application.WithConfigFile(configFile),
				application.WithModule("PubSub", pubsub.New()),
				application.WithModule("Minecraft", mcserver.New(srv)),
				application.WithModule("Command & Control", cnc.New()),
			)

			if err := app.Validate(context.Background()); err != nil {
				return err
			}

			if err := srv.ResolveVersion(context.Background()); err != nil {
				return errors.Wrapf(err, "failed to validate version %q", srv.Version)
			}

			log := logrus.WithFields(logrus.Fields{"version": srv.Version, "snapshot": srv.Snapshot, "name": srv.Name})
			log.Info("Validated server configuration")

			return nil
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "install",
		Short: "Installs the configured minecraft server",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := application.New("mcutils", "1.0.0",
				application.WithConfigFile(configFile),
				application.WithModule("PubSub", pubsub.New()),
				application.WithModule("Minecraft", mcserver.New(srv)),
				application.WithModule("Command & Control", cnc.New()),
			)

			return app.Install(context.Background())
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "run",
		Short: "Run a minecraft server",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := application.New("mcutils", "1.0.0",
				application.WithConfigFile(configFile),
				application.WithModule("PubSub", pubsub.New()),
				application.WithModule("Minecraft", mcserver.New(srv)),
				application.WithModule("Command & Control", cnc.New()),
			)

			return app.Run(context.Background())
		},
	})

	for _, cmd := range commands {
		cmd.Flags().StringVar(&srv.Version, "version", "latest", "what version to run of minecraft")
		cmd.Flags().BoolVar(&srv.Snapshot, "snapshot", false, "when version is latest, will use the latest snapshot version")
		cmd.Flags().StringVarP(&configFile, "config", "c", "", "specify an optional configuration file")
	}

	return commands
}
