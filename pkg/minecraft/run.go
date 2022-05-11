package minecraft

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/sirupsen/logrus"
)

func (s *Server) Run(ctx context.Context, log *logrus.Entry) error {
	// eula
	if err := os.WriteFile(filepath.Join(s.Path, "eula.txt"), []byte("eula=true"), 0744); err != nil {
		return errors.Wrap(err, "failed to write eula.txt")
	}

	_, jarFile := filepath.Split(s.Entrypoint())
	args := []string{
		"-Dlog4j2.formatMsgNoLookups=true",             // log4j vulnerability patching
		"-Dlog4j.configurationFile=logging-config.xml", // custom logging format so we can parse it
		fmt.Sprintf("-Xms%dg", s.InitialMemory),
		fmt.Sprintf("-Xmx%dg", s.MaxMemory),
		"-jar", jarFile,
		"nogui",
		"--nogui",
	}

	// prepend the extra arguments
	if len(s.JavaArgs) > 0 {
		args = append(s.JavaArgs, args...)
	}

	// setup our cmd
	cmd := exec.Command(java.ExecPath(s.JavaHome), args...)
	cmd.Dir, _ = filepath.Abs(filepath.FromSlash(s.Path))

	cmd.SysProcAttr = getSysProcAttr()

	stdinpipe, _ := cmd.StdinPipe()
	s.console = bufio.NewWriter(stdinpipe)

	s.fsm = s.createFSM()

	cmd.Stdout = &logParser{log: log, srv: s}
	cmd.Stderr = os.Stderr

	// output
	log.Infof("Starting server: %s", strings.Join(cmd.Args, " "))
	err := errors.Wrapf(cmd.Run(), "failed running server: %s", cmd.Path)
	_, _ = cmd.Stdout.Write([]byte("Stopped the server"))

	return err
}
