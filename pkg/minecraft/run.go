package minecraft

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/sirupsen/logrus"
)

func (s *Server) Run(ctx context.Context, log *logrus.Entry) error {
	// eula
	if err := os.WriteFile(filepath.Join(s.Path, "eula.txt"), []byte("eula=true"), 0644); err != nil {
		return errors.Wrap(err, "failed to write eula.txt")
	}

	// TODO: server.properties output

	jarpath, _ := filepath.Abs(filepath.Join(s.Path, s.VersionDetails.ID+".jar"))

	args := []string{
		"-Dlog4j2.formatMsgNoLookups=true", // log4j vulnerability patching
		"-Dlog4j.configurationFile=logging-config.xml",
		"-jar", jarpath,
		"--nogui",
	}

	cmd := exec.Command(java.ExecPath(s.JavaHome), args...)
	cmd.Dir, _ = filepath.Abs(filepath.FromSlash(s.Path))

	cmd.SysProcAttr = getSysProcAttr()

	stdinpipe, _ := cmd.StdinPipe()
	stdin := bufio.NewWriter(stdinpipe)

	cmd.Stdout = &logParser{log: log}
	cmd.Stderr = os.Stderr

	go func() {
		<-ctx.Done()
		_, _ = stdin.Write([]byte("save-all\n"))
		_ = stdin.Flush()

		_, _ = stdin.Write([]byte("stop\r\n"))
		_ = stdin.Flush()
	}()

	// output
	return errors.Wrapf(cmd.Run(), "failed running server: %s", cmd.Path)
}
