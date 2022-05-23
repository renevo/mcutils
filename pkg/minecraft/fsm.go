package minecraft

import (
	"regexp"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/looplab/fsm"
	"github.com/sirupsen/logrus"
)

const (
	StateStarting = "starting"
	StateOnline   = "online"
	StateStopping = "stopping"
	StateOffline  = "offline"
)

const (
	EventSaving            = "saving"
	EventSaved             = "saved"
	EventPlayerJoin        = "player_join"
	EventPlayerLeave       = "player_leave"
	EventWhitelistAdd      = "player_whitelist_add"
	EventWhitelistRemove   = "player_whitelist_remove"
	EventWhitelistUnknown  = "player_whitelist_unknown"
	EventWhitelistReloaded = "whitelist_reloaded"
	EventServerChat        = "server_chat"
	EventServerEmote       = "server_emote"
)

func (s *Server) createFSM() *fsm.FSM {
	return fsm.NewFSM(StateOffline, fsm.Events{
		fsm.EventDesc{Name: StateStarting, Src: []string{StateOffline}, Dst: StateStarting},
		fsm.EventDesc{Name: StateOnline, Src: []string{StateStarting}, Dst: StateOnline},
		fsm.EventDesc{Name: StateStopping, Src: []string{StateOnline, StateStarting, StateOffline}, Dst: StateStopping},
		fsm.EventDesc{Name: StateOffline, Src: []string{StateStarting, StateOnline, StateStopping}, Dst: StateOffline},
	}, fsm.Callbacks{})
}

func (s *Server) State() string {
	if s.fsm == nil {
		return StateOffline
	}

	return s.fsm.Current()
}

var stateMatchers = map[string]*regexp.Regexp{
	StateStarting: regexp.MustCompile(`^Starting minecraft server`),
	StateOnline:   regexp.MustCompile(`^Done \(\d+.\d+\w\)! For help, type \"help\"$`),
	StateStopping: regexp.MustCompile(`^Stopping the server$`),
	StateOffline:  regexp.MustCompile(`^Stopped the server$`),
}

var eventMatchers = map[string]*regexp.Regexp{
	EventSaving:            regexp.MustCompile(`^Saving the game \(this may take a moment!\)$`),
	EventSaved:             regexp.MustCompile(`Saved the game$`),
	EventWhitelistAdd:      regexp.MustCompile(`^Added (?P<player>.*) to the whitelist$`),
	EventWhitelistRemove:   regexp.MustCompile(`^Removed (?P<player>.*) from the whitelist$`),
	EventWhitelistUnknown:  regexp.MustCompile(`^That player does not exist$`),
	EventWhitelistReloaded: regexp.MustCompile(`^Reloaded the whitelist$`),
	EventServerChat:        regexp.MustCompile(`^\[Server\] (?P<msg>.*)$`),
	EventServerEmote:       regexp.MustCompile(`^\* Server (?P<msg>.*)$`),
}

func (s *Server) handleMessage(msg string, log *logrus.Entry) {
	for k, v := range stateMatchers {

		if matches := v.FindStringSubmatch(msg); len(matches) > 0 {
			if err := s.fsm.Event(k); err != nil {
				log.WithField("state", s.State()).Errorf("Failed to set state: %q", k)
			}

			log.WithField("state", s.State()).Debugf("Event: %q", k)

			if s.publisher == nil {
				continue
			}

			msg := &message.Message{
				UUID:     watermill.NewUUID(),
				Payload:  []byte(msg),
				Metadata: message.Metadata{},
			}
			for i, name := range v.SubexpNames() {
				msg.Metadata[name] = matches[i]
			}

			_ = s.publisher.Publish(k, msg)
		}
	}

	for k, v := range eventMatchers {
		if matches := v.FindStringSubmatch(msg); len(matches) > 0 {
			log.Debugf("Event: %q", k)

			if s.publisher == nil {
				continue
			}

			msg := &message.Message{
				UUID:     watermill.NewUUID(),
				Payload:  []byte(msg),
				Metadata: message.Metadata{},
			}
			for i, name := range v.SubexpNames() {
				msg.Metadata[name] = matches[i]
			}

			_ = s.publisher.Publish(k, msg)
		}
	}
}
