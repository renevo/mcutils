package minecraft

import (
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type logParser struct {
	log *logrus.Entry
	srv *Server
}

var (
	logRegex = regexp.MustCompile(`^(?:(?:\d+\-\d+-\d+\s)?\[?\d+:\d+:\d+\]?)\s\[(\w+)\]\s(.*)$`)
)

func (l *logParser) Write(d []byte) (int, error) {
	// panic prevention on logs
	defer func() {
		if r := recover(); r != nil {
			l.log.Errorf("Recovered in logParser.Write: %v", r)
		}
	}()

	lines := strings.Split(string(d), "\n")
	type entry struct {
		level   string
		message string
	}
	entries := []entry{}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		// empty lines no good
		if len(line) == 0 {
			continue
		}

		matches := logRegex.FindStringSubmatch(line)
		// didn't match the regex
		if len(matches) == 0 {
			if len(entries) == 0 {
				// no idea what it is... output it to be safe
				entries = append(entries, entry{level: "info", message: line})
				continue
			}

			// append to the previous log entry
			entries[len(entries)-1].message = entries[len(entries)-1].message + "\n" + raw
			continue
		}

		// didn't match the regex correctly?
		if len(matches) != 3 {
			entries = append(entries, entry{level: "info", message: line})
			continue
		}

		entries = append(entries, entry{level: strings.ToLower(matches[1]), message: matches[2]})
	}

	for _, entry := range entries {
		msg := strings.TrimSpace(entry.message)
		l.srv.handleMessage(msg, l.log)

		stateLogger := l.log.WithField("state", l.srv.State())

		outputFn := stateLogger.Error // don't know what else there might be?

		// level switch
		switch entry.level {
		case "trace":
			outputFn = stateLogger.Trace

		case "debug":
			outputFn = stateLogger.Debug

		case "info":
			outputFn = stateLogger.Info

		case "warn":
			outputFn = stateLogger.Warning

		}

		outputFn(msg)

	}

	return len(d), nil
}
