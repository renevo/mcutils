package discord

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/renevo/mcutils/pkg/minecraft"
)

func (m *module) handleEvents(ctx context.Context, ch <-chan *message.Message) {
	for msg := range ch {
		msg.Ack()

		switch msg.Metadata["event"] {
		case minecraft.EventPlayerJoin:
			player, ok := msg.Metadata["player"]
			if !ok {
				continue
			}

			m.sendMessage("> **%s** *joined the game*", player)

		case minecraft.EventPlayerLeave:
			player, ok := msg.Metadata["player"]
			if !ok {
				continue
			}

			m.sendMessage("> **%s** *left the game*", player)

		case minecraft.StateStarting:
			_ = m.session.UpdateGameStatus(0, "Minecraft Server Startup")

		case minecraft.StateOnline:
			_ = m.session.UpdateGameStatus(0, "Minecraft Server")

		case minecraft.StateStopping:
			_ = m.session.UpdateGameStatus(0, "Minecraft Server Stopping")

		case minecraft.StateOffline:
			_ = m.session.UpdateGameStatus(0, "")

		case minecraft.EventChat:
			player, ok := msg.Metadata["player"]
			if !ok {
				continue
			}
			message, ok := msg.Metadata["message"]
			if !ok {
				continue
			}

			// don't echo server stuff
			if player == "Server" {
				continue
			}

			m.sendMessage("> [**%s**]: %s", player, message)

		case minecraft.EventPlayerAdvancement:
			player, ok := msg.Metadata["player"]
			if !ok {
				continue
			}
			advancement, ok := msg.Metadata["advancement"]
			if !ok {
				continue
			}

			m.sendMessage("> **%s** *got advancement* **%s**", player, advancement)
		}
	}
}
