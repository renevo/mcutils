package server

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/internal/download"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/renevo/mcutils/pkg/minecraft/version"
)

func (s *Server) Install(ctx context.Context) (version.Version, error) {
	// get the manifest
	manifest, err := version.GetManifest(context.Background())
	if err != nil {
		return s.version, err
	}

	// looku up the version
	lookupVersion := s.Version
	if s.Version == "latest" && s.Snapshot {
		lookupVersion = "snapshot"
	}

	v, err := manifest.GetVersion(context.Background(), lookupVersion)
	if err != nil {
		return s.version, errors.Wrapf(err, "failed to get version %q", s.Version)
	}

	s.version = v

	// create the output directory
	if err := os.MkdirAll(s.Path, 0644); err != nil {
		return s.version, errors.Wrapf(err, "failed to create server directory %q", s.Path)
	}

	// only download if not already downloaded
	jarPath := s.Entrypoint()
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		if err := download.File(ctx, v.Downloads.Server.URL, jarPath); err != nil {
			return s.version, errors.Wrapf(err, "failed to download %q", v.Downloads.Server.URL)
		}
	}

	// setup the java home
	if s.JavaHome == "" {
		s.JavaHome, _ = filepath.Abs(filepath.Join(s.Path, java.VersionPaths[s.version.Java.Version]))
	}

	// if the java home doesn't exist, try to install the correct java version
	if _, err := os.Stat(s.JavaHome); os.IsNotExist(err) {
		if err := java.Versions.Install(ctx, s.version.Java.Version, s.Path); err != nil {
			return s.version, errors.Wrapf(err, "failed to download and install java to %q", s.JavaHome)
		}
	}

	if err := java.WriteLoggingConfig(s.Path); err != nil {
		return s.version, err
	}

	return s.version, nil
}
