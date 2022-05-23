package minecraft

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/renevo/rpc"
	"github.com/spf13/cobra"
)

func clientCommands() []*cobra.Command {
	address := "127.0.0.1:2311"
	token := ""

	commands := []*cobra.Command{}

	rpcFunc := func(method string, command string, reply any) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		client, err := rpc.Dial(ctx, "tcp", address)
		if err != nil {
			return errors.Wrapf(err, "failed to dial %q", address)
		}

		ctx = rpc.ContextWithHeaders(ctx, rpc.Header{}.Set(rpcHeaderToken, token))

		if err := client.Call(ctx, "Minecraft."+method, command, reply); err != nil {
			return errors.Wrapf(err, "failed to execute %q", command)
		}

		return nil
	}

	commands = append(commands, &cobra.Command{
		Use:   "exec",
		Short: "Execute a command against a remote server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reply := false
			return rpcFunc("Execute", args[0], &reply)
		},
	})

	whitelistCommand := &cobra.Command{
		Use:   "whitelist",
		Short: "Whitelist commands",
	}

	whitelistCommand.AddCommand(&cobra.Command{
		Use:   "add",
		Short: "Add a player from the whitelist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			playerActual := ""
			if err := rpcFunc("WhitelistAdd", args[0], &playerActual); err != nil {
				return err
			}

			fmt.Printf("Player %q added to the whitelist\n", playerActual)
			return nil
		},
	})

	whitelistCommand.AddCommand(&cobra.Command{
		Use:   "remove",
		Short: "Remove a player from the whitelist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			playerActual := ""
			if err := rpcFunc("WhitelistRemove", args[0], &playerActual); err != nil {
				return err
			}

			fmt.Printf("Player %q removed from the whitelist\n", playerActual)
			return nil
		},
	})

	commands = append(commands, whitelistCommand)

	for _, cmd := range commands {
		cmd.PersistentFlags().StringVarP(&address, "address", "a", address, "Specify the address:port for the server rpc endpoint")
		cmd.PersistentFlags().StringVarP(&token, "token", "t", token, "token to use when making requests to the server")
	}

	return commands
}
