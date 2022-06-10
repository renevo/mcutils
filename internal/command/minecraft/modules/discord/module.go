package discord

import (
	"context"

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

	registeredCommands []*discordgo.ApplicationCommand
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

	if err := m.addCommands(ctx, m.session); err != nil {
		return errors.Wrap(err, "failed to create discord commands")
	}

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
		m.removeCommands(ctx, m.session)
		_ = m.session.Close()
	}

	// stop our pub/sub
	if m.cancel != nil {
		m.cancel()
	}

	return nil
}
