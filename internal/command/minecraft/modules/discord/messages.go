package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
)

func (m *module) handleMessages(ctx context.Context) func(*discordgo.Session, *discordgo.MessageCreate) {
	srv := ext.Minecraft(ctx)

	return func(s *discordgo.Session, msg *discordgo.MessageCreate) {
		// only specified channel
		if msg.ChannelID != m.cfg.ChannelID {
			return
		}

		// don't get into an echo loop....
		if msg.Author.ID == s.State.User.ID {
			return
		}

		// render the contents away from @<id> things to @User things
		content := msg.ContentWithMentionsReplaced()

		// some things can be just images, etc... lets not spam with nothingness
		if len(content) == 0 {
			return
		}

		// use tellraw so it doesn't look like it came from the server
		// content may have issues, and emojii don't work (show as [] character)
		_ = srv.ExecuteCommand(fmt.Sprintf(`tellraw @a {"text": %q, "color": "blue", "extra": [ { "text": %q, "color": "white" } ] }`,
			"<"+msg.Author.Username+">",
			" "+content,
		))
	}
}
