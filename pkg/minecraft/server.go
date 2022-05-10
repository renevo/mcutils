package minecraft

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/renevo/mcutils/pkg/minecraft/version"
)

// Server for Minecraft
type Server struct {
	Name       string            `hcl:"name,label" property:"motd"`
	Path       string            `hcl:"path"`
	Version    string            `hcl:"version"`
	Snapshot   bool              `hcl:"snapshot,optional"`
	JavaHome   string            `hcl:"java_home,optional"`
	Properties map[string]string `hcl:"properties,optional"`

	VersionDetails version.Version
}

// Default will return a default configured Minecraft server
func Default() *Server {
	return &Server{
		Name:       "minecraft",
		Path:       "./.minecraft/",
		Version:    "latest",
		JavaHome:   os.Getenv("JAVA_HOME"),
		Properties: make(map[string]string),
	}
}

// Entrypoint returns the location of the Minecraft server jar
func (s *Server) Entrypoint() string {
	if s.VersionDetails.ID == "" {
		return filepath.Join(s.Path, "server.jar")
	}

	return filepath.Join(s.Path, s.VersionDetails.ID+".jar")
}

func (s *Server) ResolveVersion(ctx context.Context) error {
	// get the manifest
	manifest, err := version.GetManifest(ctx)
	if err != nil {
		return err
	}

	// lookup the version
	lookupVersion := s.Version
	if s.Version == "latest" && s.Snapshot {
		lookupVersion = "snapshot"
	}

	v, err := manifest.GetVersion(ctx, lookupVersion)
	if err != nil {
		return errors.Wrapf(err, "failed to get version %q", s.Version)
	}

	s.VersionDetails = v

	// setup the java home
	if s.JavaHome == "" {
		s.JavaHome, _ = filepath.Abs(filepath.Join(s.Path, java.VersionPaths[s.VersionDetails.Java.Version]))
	}

	return nil
}
