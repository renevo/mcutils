package minecraft

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/internal/download"
	"github.com/renevo/mcutils/pkg/java"
	"github.com/renevo/mcutils/pkg/minecraft/version"
)

func (s *Server) Install(ctx context.Context) (version.Version, error) {
	if s.VersionDetails.ID == "" {
		if err := s.ResolveVersion(ctx); err != nil {
			return s.VersionDetails, err
		}
	}

	// create the output directory
	if err := os.MkdirAll(s.Path, 0644); err != nil {
		return s.VersionDetails, errors.Wrapf(err, "failed to create server directory %q", s.Path)
	}

	// only download if not already downloaded
	jarPath := s.Entrypoint()
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		if err := download.File(ctx, s.VersionDetails.Downloads.Server.URL, jarPath); err != nil {
			return s.VersionDetails, errors.Wrapf(err, "failed to download %q", s.VersionDetails.Downloads.Server.URL)
		}
	}

	// if the java home doesn't exist, try to install the correct java version
	if _, err := os.Stat(s.JavaHome); os.IsNotExist(err) {
		if err := java.Versions.Install(ctx, s.VersionDetails.Java.Version, s.Path); err != nil {
			return s.VersionDetails, errors.Wrapf(err, "failed to download and install java to %q", s.JavaHome)
		}
	}

	if err := java.WriteLoggingConfig(s.Path); err != nil {
		return s.VersionDetails, err
	}

	return s.VersionDetails, nil
}
