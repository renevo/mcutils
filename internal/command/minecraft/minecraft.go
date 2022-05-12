package minecraft

import "github.com/spf13/cobra"

const (
	rpcHeaderToken = "X-Minecraft-Token"
)

func New() *cobra.Command {
	minecraftCommand := &cobra.Command{
		Use:   "minecraft",
		Short: "Interact with a minecraft server",
	}

	minecraftCommand.AddCommand(clientCommands()...)
	minecraftCommand.AddCommand(serverCommands()...)

	return minecraftCommand
}
