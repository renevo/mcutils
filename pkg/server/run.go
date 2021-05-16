package server

import (
	"context"
	"fmt"
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
		"-Dlog4j.configurationFile=logging-config.xml",
		"-jar", jarpath,
	}

	cmd := exec.Command(java.ExecPath(s.JavaHome), args...)
	cmd.Dir = s.Path

	// build up the environment
	cmd.Env = []string{}
	cmd.Env = append(cmd.Env, fmt.Sprintf("JAVA_HOME=%s", s.JavaHome))
	cmd.Env = append(cmd.Env, "HOME", s.Path)
	cmd.Env = append(cmd.Env, "APPDATA", s.Path)

	// temporary....
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to run server")
	}

	return nil
}
