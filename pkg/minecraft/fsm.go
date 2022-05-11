package minecraft

import (
	"regexp"

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
	EventSaving      = "saving"
	EventSaved       = "saved"
	EventPlayerJoin  = "player_join"
	EventPlayerLeave = "player_leave"
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
	EventSaving: regexp.MustCompile(`^Saving the game \(this may take a moment!\)$`),
	EventSaved:  regexp.MustCompile(`Saved the game$`),
}

func (s *Server) handleMessage(msg string, log *logrus.Entry) {
	for k, v := range stateMatchers {
		if v.MatchString(msg) {
			if err := s.fsm.Event(k); err != nil {
				log.WithField("state", s.State()).Errorf("Failed to set state: %q", k)
			}

			// TODO: publish event for k
			log.WithField("state", s.State()).Infof("Event: %q", k)
		}
	}

	for k, v := range eventMatchers {
		if v.MatchString(msg) {
			// TODO: publish event for k
			log.Infof("Event: %q", k)
		}
	}
}
