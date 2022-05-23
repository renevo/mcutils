package minecraft

import (
	"context"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/pkg/errors"
	"github.com/portcullis/application"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/cnc"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/mcserver"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/pubsub"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/renevo/mcutils/pkg/minecraft"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func serverCommands() []*cobra.Command {
	configFile := "./minecraft.hcl"

	srv := minecraft.Default()
	minecraftConfig := struct {
		Minecraft []*minecraft.Server `hcl:"minecraft,block"`
	}{
		Minecraft: []*minecraft.Server{srv},
	}

	loadConfig := func() error {
		if len(configFile) > 0 {
			if err := hclsimple.DecodeFile(configFile, nil, &minecraftConfig); err != nil {
				return errors.Wrap(err, "failed to parse config file")
			}

			if len(minecraftConfig.Minecraft) != 1 {
				return errors.New("you must specify exactly one server block in the configuration file")
			}
		}

		return nil
	}

	commands := []*cobra.Command{}

	commands = append(commands, &cobra.Command{
		Use:   "validate",
		Short: "Validates configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}

			if err := srv.ResolveVersion(context.Background()); err != nil {
				return errors.Wrapf(err, "failed to validate version %q", srv.Version)
			}

			log := logrus.WithFields(logrus.Fields{"version": srv.Version, "snapshot": srv.Snapshot, "name": srv.Name})
			log.Info("Validated server configuration")
			log.Printf("%+v", srv)

			return nil
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "install",
		Short: "Installs the configured minecraft server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}

			log := logrus.WithFields(logrus.Fields{"version": srv.Version, "snapshot": srv.Snapshot, "name": srv.Name})
			log.Info("Installing server and dependencies")

			v, err := srv.Install(context.Background())
			if err != nil {
				return errors.Wrap(err, "failed to install server")
			}

			log.Infof("Version ID: %q; Type: %q; URL: %q;", v.ID, v.Type, v.Downloads.Server.URL)
			log.Infof("JAVA_HOME: %q", srv.JavaHome)
			log.Infof("Exec Path: %q", java.ExecPath(srv.JavaHome))

			return nil
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "run",
		Short: "Run a minecraft server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return application.Run("mcutils", "1.0.0",
				application.WithConfigFile(configFile),
				application.WithModule("PubSub", pubsub.New()),
				application.WithModule("Minecraft", mcserver.New()),
				application.WithModule("Command & Control", cnc.New()),
			)
		},
	})

	for _, cmd := range commands {
		cmd.Flags().StringVar(&srv.Version, "version", "latest", "what version to run of minecraft")
		cmd.Flags().BoolVar(&srv.Snapshot, "snapshot", false, "when version is latest, will use the latest snapshot version")
		cmd.Flags().StringVarP(&configFile, "config", "c", "", "specify an optional configuration file")
	}

	return commands
}
