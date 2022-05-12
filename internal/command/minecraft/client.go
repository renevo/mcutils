package minecraft

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/renevo/rpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func clientCommands() []*cobra.Command {
	address := "127.0.0.1:2311"
	token := ""

	commands := []*cobra.Command{}

	rpcFunc := func(method string, command string) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		client, err := rpc.Dial(ctx, "tcp", address)
		if err != nil {
			return errors.Wrapf(err, "failed to dial %q", address)
		}

		ctx = rpc.ContextWithHeaders(ctx, rpc.Header{}.Set(rpcHeaderToken, token))

		reply := ""
		if err := client.Call(ctx, "Minecraft."+method, command, &reply); err != nil {
			return errors.Wrapf(err, "failed to execute %q", command)
		}

		logrus.Infof("Executed %q with reply %q", command, reply)
		return nil
	}

	commands = append(commands, &cobra.Command{
		Use:   "exec",
		Short: "Execute a command against a remote server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rpcFunc("Execute", args[0])
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "whitelist",
		Short: "Whitelist a user a command against a remote server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rpcFunc("WhitelistAdd", args[0])
		},
	})

	for _, cmd := range commands {
		cmd.Flags().StringVarP(&address, "address", "a", address, "Specify the address:port for the server rpc endpoint")
		cmd.Flags().StringVarP(&token, "token", "t", token, "token to use when making requests to the server")
	}

	return commands
}
