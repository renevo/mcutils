package minecraft

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/pkg/errors"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/renevo/mcutils/pkg/minecraft"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	configFile := "./server.hcl"

	srv := minecraft.Default()
	serverConfig := struct {
		Server []*minecraft.Server `hcl:"server,block"`
	}{
		Server: []*minecraft.Server{srv},
	}

	loadConfig := func() error {
		if len(configFile) > 0 {
			if err := hclsimple.DecodeFile(configFile, nil, &serverConfig); err != nil {
				return errors.Wrap(err, "failed to parse config file")
			}

			if len(serverConfig.Server) != 1 {
				return errors.New("you must specify exactly one server block in the configuration file")
			}
		}

		return nil
	}

	minecraftCommand := &cobra.Command{
		Use:   "minecraft",
		Short: "Runs a minecraft server",
		Long:  `Will download, configure, and run a Minecraft server. Without any configuration, this will be a default vanilla minecraft server on localhost`,
	}

	minecraftCommand.AddCommand(&cobra.Command{
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

	minecraftCommand.AddCommand(&cobra.Command{
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

	minecraftCommand.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Run a minecraft server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}

			log := logrus.WithFields(logrus.Fields{"version": srv.Version, "snapshot": srv.Snapshot, "name": srv.Name})
			log.Info("Installing server and dependencies")

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			v, err := srv.Install(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to install server")
			}

			sigCh := make(chan os.Signal, 2)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			go func() {
				sig := <-sigCh
				log.Infof("Stopping server... %v", sig)
				cancel()
			}()

			log.Infof("Version ID: %q; Type: %q; URL: %q;", v.ID, v.Type, v.Downloads.Server.URL)
			log.Infof("JAVA_HOME: %q", srv.JavaHome)
			log.Infof("Exec Path: %q", java.ExecPath(srv.JavaHome))

			err = srv.Run(ctx, log)

			if err != nil {
				log.Infof("Stopped Server: %v", err)
			} else {
				log.Infof("Stopped Server")
			}

			return err
		},
	})

	minecraftCommand.PersistentFlags().StringVar(&srv.Version, "version", "latest", "what version to run of minecraft")
	minecraftCommand.PersistentFlags().BoolVar(&srv.Snapshot, "snapshot", false, "when version is latest, will use the latest snapshot version")
	minecraftCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "specify an optional configuration file")

	return minecraftCommand
}
