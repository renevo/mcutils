package server

import (
	"os"
	"path/filepath"

	"github.com/renevo/mcutils/pkg/minecraft/version"
)

// Server for Minecraft
type Server struct {
	Name     string `hcl:"name,label"`
	Path     string `hcl:"path"`
	Version  string `hcl:"version"`
	Snapshot bool   `hcl:"snapshot"`
	JavaHome string `hcl:"java_home"`

	version version.Version
}

// Default will return a default configured Minecraft server
func Default() *Server {
	return &Server{
		Name:     "minecraft",
		Path:     "./.minecraft/",
		Version:  "latest",
		Snapshot: false,
		JavaHome: os.Getenv("JAVA_HOME"),
	}
}

// Entrypoint returns the location of the Minecraft server jar
func (s *Server) Entrypoint() string {
	if s.version.ID == "" {
		return filepath.Join(s.Path, "server.jar")
	}

	return filepath.Join(s.Path, s.version.ID+".jar")
}
