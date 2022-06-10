package discord

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
)

func (m *module) addCommands(ctx context.Context, s *discordgo.Session) error {
	srv := ext.Minecraft(ctx)

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "minecraft",
			Description: "Minecraft Commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "status",
					Description: "Server Status",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
	}

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"minecraft": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			if len(options) == 0 {
				content := "Unknown command"
				_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
			}

			switch options[0].Name {
			case "status":

				status, err := srv.Status()
				if err != nil {
					content := "Status not enabled"
					_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: content,
						},
					})

					return
				}

				// really basic for now
				content := fmt.Sprintf(`Minecraft Server: *%s*
				
**Map**: %s
**Version**: %s
**Players**: %d/%d
%s`,
					status.Description,
					status.MapName,
					status.Version.Name,
					status.Players.Online,
					status.Players.Max,
					strings.Join(status.Players.PlayerList, "\n"))

				_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})

			default:
				content := "Unknown command"
				_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
			}
		},
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	m.registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, m.cfg.ServerID, v)
		if err != nil {
			return errors.Wrapf(err, "failed to create command %q", v.Name)
		}

		m.registeredCommands[i] = cmd
	}

	return nil
}

func (m *module) removeCommands(ctx context.Context, s *discordgo.Session) {
	for _, v := range m.registeredCommands {
		_ = s.ApplicationCommandDelete(s.State.User.ID, m.cfg.ServerID, v.ID)
	}
}
