package discord

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/portcullis/application"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
	"github.com/renevo/mcutils/pkg/minecraft"
)

type cfg struct {
	ServerID    string `hcl:"discord_server_id,optional"`
	ChannelID   string `hcl:"discord_server_channel,optional"`
	ServerToken string `hcl:"discord_server_token,optional" env:"DISCORD_SERVER_TOKEN"`
}

type module struct {
	cfg     *cfg
	session *discordgo.Session
	ctx     context.Context
	cancel  context.CancelFunc
}

func New() application.Module {
	return &module{
		cfg: &cfg{},
	}
}

func (m *module) Config() (interface{}, error) {
	return m.cfg, nil
}

func (m *module) Start(ctx context.Context) error {
	if m.cfg.ChannelID == "" || m.cfg.ServerID == "" || m.cfg.ServerToken == "" {
		return nil
	}

	log := ext.Logger(ctx)

	session, err := discordgo.New("Bot " + m.cfg.ServerToken)
	if err != nil {
		return errors.Wrap(err, "failed to create discord session")
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	session.AddHandler(m.handleMessages(ctx))
	session.Identify.Intents = discordgo.IntentGuildMessages

	if err := session.Open(); err != nil {
		return errors.Wrap(err, "failed to connect to discord")
	}

	m.session = session

	subscriber := ext.Subscriber(ctx)
	m.ctx, m.cancel = context.WithCancel(context.Background())
	ch, err := subscriber.Subscribe(m.ctx, minecraft.EventAll)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to server state")
	}

	go m.handleEvents(ctx, ch)

	return nil
}

func (m *module) Stop(ctx context.Context) error {
	if m.session != nil {
		_ = m.session.Close()
	}

	// stop our pub/sub
	if m.cancel != nil {
		m.cancel()
	}

	return nil
}

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

		content := msg.ContentWithMentionsReplaced()
		if len(content) == 0 {
			return
		}

		// use tellraw so it doesn't look like it came from the server
		_ = srv.ExecuteCommand(fmt.Sprintf(`tellraw @a {"text": %q, "color": "blue", "extra": [ { "text": %q, "color": "white" } ] }`,
			"<"+msg.Author.Username+">",
			" "+content,
		))
	}
}

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

func (m *module) sendMessage(msg string, args ...any) {
	if m.session == nil {
		return
	}

	if _, err := m.session.ChannelMessageSend(m.cfg.ChannelID, fmt.Sprintf(msg, args...)); err != nil {
		fmt.Println(err.Error())
	}
}
