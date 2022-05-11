package minecraft

import "github.com/spf13/cobra"

func New() *cobra.Command {
	minecraftCommand := &cobra.Command{
		Use:   "minecraft",
		Short: "Interact with a minecraft server",
	}

	minecraftCommand.AddCommand(clientCommands()...)
	minecraftCommand.AddCommand(serverCommands()...)

	return minecraftCommand
}
