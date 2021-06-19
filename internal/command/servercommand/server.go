package servercommand

import (
	"context"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/pkg/errors"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/renevo/mcutils/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	configFile := "./server.hcl"

	srv := server.Default()
	serverConfig := struct {
		Server []*server.Server `hcl:"server,block"`
	}{
		Server: []*server.Server{srv},
	}

	serverCommand := &cobra.Command{
		Use:   "server",
		Short: "Runs a minecraft server",
		Long:  `Will download, configure, and run a Minecraft server. Without any configuration, this will be a default vanilla minecraft server on localhost`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(configFile) > 0 {
				if err := hclsimple.DecodeFile(configFile, nil, &serverConfig); err != nil {
					return errors.Wrap(err, "failed to parse config file")
				}

				if len(serverConfig.Server) != 1 {
					return errors.New("you must specify exactly one server block in the configuration file")
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logrus.WithFields(logrus.Fields{"version": srv.Version, "snapshot": srv.Snapshot, "name": srv.Name})
			log.Info("Initialize the server here....")

			v, err := srv.Install(context.Background())
			if err != nil {
				return errors.Wrap(err, "failed to install server")
			}

			log.Infof("Version ID: %q; Type: %q; URL: %q;", v.ID, v.Type, v.Downloads.Server.URL)
			log.Infof("JAVA_HOME: %q", srv.JavaHome)
			log.Infof("Exec Path: %q", java.ExecPath(srv.JavaHome))

			// will do more later...
			return errors.Wrapf(srv.Run(context.Background()), "failed to run server")
		},
	}

	serverCommand.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validates configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logrus.WithFields(logrus.Fields{"version": srv.Version, "snapshot": srv.Snapshot, "name": srv.Name})
			log.Info("Validated server configuration")
			log.Printf("%+v", srv)

			return nil
		},
	})

	serverCommand.PersistentFlags().StringVar(&srv.Version, "version", "latest", "what version to run of minecraft")
	serverCommand.PersistentFlags().BoolVar(&srv.Snapshot, "snapshot", false, "when version is latest, will use the latest snapshot version")
	serverCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "specify an optional configuration file")

	return serverCommand
}
