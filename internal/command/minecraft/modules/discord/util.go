package discord

import "fmt"

func (m *module) sendMessage(msg string, args ...any) {
	if m.session == nil {
		return
	}

	_, _ = m.session.ChannelMessageSend(m.cfg.ChannelID, fmt.Sprintf(msg, args...))
}
