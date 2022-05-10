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
	logRegex = regexp.MustCompile(`^(?:\d+\-\d+-\d+\s\d+:\d+:\d+)\s\[(\w+)\]\s(.*)$`)
)

func (l *logParser) Write(d []byte) (int, error) {
	lines := strings.Split(string(d), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		matches := logRegex.FindStringSubmatch(line)
		if len(matches) == 0 {
			//l.log.Info(line)
			continue
		}

		if len(matches) != 3 {
			l.log.Info(matches[0])
			continue
		}

		outputFn := l.log.Print

		// level
		switch strings.ToLower(matches[1]) {
		case "info":
			outputFn = l.log.Info
		case "warn":
			outputFn = l.log.Warning
		case "error":
			outputFn = l.log.Error
		}

		outputFn(matches[2])
	}
	return len(d), nil
}
