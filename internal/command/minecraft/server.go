package minecraft

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/pkg/errors"
	"github.com/renevo/mcutils/internal/control"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/renevo/mcutils/pkg/minecraft"
	"github.com/renevo/rpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func serverCommands() []*cobra.Command {
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
			if err := loadConfig(); err != nil {
				return err
			}

			log := logrus.WithFields(logrus.Fields{"version": srv.Version, "snapshot": srv.Snapshot, "name": srv.Name})
			if srv.FabricJar() != "" {
				log = log.WithFields(logrus.Fields{"flavor": "fabric", "fabric": srv.FabricVersionLoader})
			} else {
				log = log.WithField("flavor", "vanilla")
			}

			log.Info("Installing server and dependencies")

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			v, err := srv.Install(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to install server")
			}

			sigCh := make(chan os.Signal, 2)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			log.Infof("Version ID: %q; Type: %q; URL: %q;", v.ID, v.Type, v.Downloads.Server.URL)
			log.Infof("JAVA_HOME: %q", srv.JavaHome)
			log.Infof("Exec Path: %q", java.ExecPath(srv.JavaHome))

			// TODO: add support for sigint == save-all
			// TODO: add shutdown timeout to just kill the java pid
			go func() {
				sig := <-sigCh
				log.Infof("Stopping server... %v", sig)
				if err := srv.ExecuteCommand("save-all"); err != nil {
					log.Errorf("Failed to save: %v", err)
				}
				if err := srv.ExecuteCommand("stop"); err != nil {
					log.Errorf("Failed to stop - server may be zombied: %v", err)
				}

				cancel()
			}()

			if srv.ControlAddr != "" {
				controlServer := rpc.NewServer()
				if err := controlServer.Register(&control.Minecraft{Server: srv}); err != nil {
					panic(err)
				}

				ln, err := net.Listen("tcp", srv.ControlAddr)
				if err != nil {
					return errors.Wrapf(err, "failed to open control server: %q", srv.ControlAddr)
				}
				defer ln.Close()

				log.Infof("RPC Control Server Running on %v", ln.Addr().String())

				// TODO: probably implement a token ... so its not wide open to the world to do stuff
				controlServer.Use(func(next rpc.MiddlewareHandler) rpc.MiddlewareHandler {
					return func(ctx context.Context, rw rpc.ResponseWriter, req *rpc.Request) {
						start := time.Now()
						log.Infof("RPC executing %q", req.ServiceMethod)
						next(ctx, rw, req)
						log.Infof("RPC executed %q in %v", req.ServiceMethod, time.Since(start))
					}
				})

				go func() {
					controlServer.Accept(ctx, ln)
				}()
			}

			err = srv.Run(ctx, log)

			if err != nil {
				log.Errorf("Stopped Server: %v", err)
			} else {
				log.Infof("Stopped Server")
			}

			return err
		},
	})

	for _, cmd := range commands {
		cmd.Flags().StringVar(&srv.Version, "version", "latest", "what version to run of minecraft")
		cmd.Flags().BoolVar(&srv.Snapshot, "snapshot", false, "when version is latest, will use the latest snapshot version")
		cmd.Flags().StringVarP(&configFile, "config", "c", "", "specify an optional configuration file")
	}

	return commands
}
