package java

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func WriteLoggingConfig(path string) error {
	consoleConfig := `<?xml version="1.0" encoding="UTF-8"?>
	<Configuration>
		<Appenders>
			<Console name="console" target="SYSTEM_OUT">
				<PatternLayout pattern="%d{yyyy-MM-dd HH:mm:ss} [%level] %msg{nolookups}%n" />
			</Console>
		</Appenders>
		<Loggers>
			<Root level="debug">
				<AppenderRef ref="console" />
			</Root>
		</Loggers>
	</Configuration>`

	outputFile := filepath.Join(path, "logging-config.xml")

	if err := os.WriteFile(outputFile, []byte(consoleConfig), 0744); err != nil {
		return errors.Wrapf(err, "failed to write logging configuration %q", outputFile)
	}

	return nil
}
