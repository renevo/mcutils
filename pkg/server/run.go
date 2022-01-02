package server

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/pkg/java"
)

func (s *Server) Run(ctx context.Context) error {
	// eula
	if err := os.WriteFile(filepath.Join(s.Path, "eula.txt"), []byte("eula=true"), 0644); err != nil {
		return errors.Wrap(err, "failed to write eula.txt")
	}

	// TODO: server.properties output

	jarpath, _ := filepath.Abs(filepath.Join(s.Path, s.version.ID+".jar"))

	args := []string{
		"-Dlog4j2.formatMsgNoLookups=true", // log4j vulnerability patching
		"-Dlog4j.configurationFile=logging-config.xml",
		"-jar", jarpath,
		"--nogui",
	}

	cmd := exec.Command(java.ExecPath(s.JavaHome), args...)
	cmd.Dir, _ = filepath.Abs(filepath.FromSlash(s.Path))

	// on windows we can't set env to clear then run, as it seems that it "needs" some stuff....

	// temporary....
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// idealistically we would run an background go routine here
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to run server")
	}

	return nil
}
