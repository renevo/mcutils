package minecraft

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/renevo/mcutils/pkg/minecraft/version"
)

// Server for Minecraft
type Server struct {
	Name           string     `hcl:"name,label" property:"motd"`
	Path           string     `hcl:"path"`
	Version        string     `hcl:"version"`
	Snapshot       bool       `hcl:"snapshot,optional"`
	JavaHome       string     `hcl:"java_home,optional"`
	InitialMemory  int        `hcl:"memory_min,optional"`
	MaxMemory      int        `hcl:"memory_max,optional"`
	JavaArgs       []string   `hcl:"java_extra_args,optional"`
	Properties     Properties `hcl:"properties,optional"`
	VersionDetails version.Version

	FabricVersionLoader    string `hcl:"fabric_loader,optional"`
	FabricVersionInstaller string `hcl:"fabric_installer,optional"`

	console *bufio.Writer
}

// Default will return a default configured Minecraft server
func Default() *Server {
	return &Server{
		Name:          "minecraft",
		Path:          "./.minecraft/",
		Version:       "latest",
		JavaHome:      os.Getenv("JAVA_HOME"),
		Properties:    Properties{},
		InitialMemory: 1,
		MaxMemory:     2,
	}
}

// Entrypoint returns the location of the Minecraft server jar
func (s *Server) Entrypoint() string {
	if fabric := s.FabricJar(); fabric != "" {
		return fabric
	}

	return s.MinecraftJar()
}

// FabricJar to use
func (s *Server) FabricJar() string {
	if s.FabricVersionInstaller != "" && s.FabricVersionLoader != "" {
		return filepath.Join(s.Path, fmt.Sprintf("fabric-server-mc.%s-loader.%s-launcher.%s.jar", s.VersionDetails.ID, s.FabricVersionLoader, s.FabricVersionInstaller))
	}

	return ""
}

// Minecraft jar to use
func (s *Server) MinecraftJar() string {
	if s.VersionDetails.ID == "" {
		return filepath.Join(s.Path, "server.jar")
	}

	return filepath.Join(s.Path, s.VersionDetails.ID+".jar")
}

// ResolveVersion from the minecraft version manifect API
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

// ExecuteCommand against the server, this is a standard minecraft command
func (s *Server) ExecuteCommand(cmd string) error {
	if s.console == nil {
		return errors.New("server console not available")
	}

	if _, err := s.console.Write([]byte(cmd + "\n")); err != nil {
		return errors.Wrapf(err, "failed to send command %q", cmd)
	}

	return errors.Wrapf(s.console.Flush(), "failed to flush command %q to console", cmd)
}
