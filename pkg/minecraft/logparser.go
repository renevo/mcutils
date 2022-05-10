package minecraft

import (
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type logParser struct {
	log *logrus.Entry
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
		if len(line) == 0 {
			continue
		}

		matches := logRegex.FindStringSubmatch(line)
		if len(matches) == 0 {
			if len(entries) == 0 {
				// no idea what it is... output it to be safe
				l.log.Info(line)
				continue
			}

			// append to the previous log entry
			entries[len(entries)-1].message = entries[len(entries)-1].message + "\n" + raw
			continue
		}

		if len(matches) != 3 {
			entries = append(entries, entry{level: "info", message: line})
			continue
		}

		entries = append(entries, entry{level: strings.ToLower(matches[1]), message: matches[2]})
	}

	for _, entry := range entries {
		outputFn := l.log.Fatal // don't know what else there might be?

		// level switch
		switch entry.level {
		case "trace":
			outputFn = l.log.Trace

		case "debug":
			outputFn = l.log.Debug

		case "info":
			outputFn = l.log.Info

		case "warn":
			outputFn = l.log.Warning

		case "error":
			outputFn = l.log.Error
		}

		outputFn(strings.TrimSpace(entry.message))
	}

	return len(d), nil
}
