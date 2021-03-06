package minecraft

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	if err := os.MkdirAll(s.Path, 0744); err != nil {
		return s.VersionDetails, errors.Wrapf(err, "failed to create server directory %q", s.Path)
	}

	// only download if not already downloaded
	jarPath := s.MinecraftJar()
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

	// server.properties output
	if err := s.WriteProperties(); err != nil {
		return s.VersionDetails, errors.Wrap(err, "failed to merge server properties")
	}

	worldName := s.Properties["level-name"]
	if worldName == "" {
		worldName = "world"
	}

	// datapacks
	if s.PurgeDatapacks {
		if err := os.RemoveAll(filepath.Join(s.Path, worldName, "datapacks")); err != nil {
			return s.VersionDetails, errors.Wrap(err, "failed to purge datapacks directory")
		}
	}

	if err := s.Datapacks.Install(ctx, filepath.Join(s.Path, worldName)); err != nil {
		return s.VersionDetails, errors.Wrap(err, "failed to install datapacks")
	}

	// optional fabric
	if s.FabricVersionInstaller != "" && s.FabricVersionLoader != "" {
		fabricFile := s.FabricJar()
		fabricSettings := filepath.Join(s.Path, "fabric-server-launcher.properties")

		if _, err := os.Stat(fabricFile); os.IsNotExist(err) {
			if err := download.File(ctx, fmt.Sprintf("https://meta.fabricmc.net/v2/versions/loader/%s/%s/%s/server/jar", s.VersionDetails.ID, s.FabricVersionLoader, s.FabricVersionInstaller), fabricFile); err != nil {
				return s.VersionDetails, errors.Wrapf(err, "failed to install fabric %q", fabricFile)
			}
		}

		if err := os.WriteFile(fabricSettings, []byte(fmt.Sprintf("serverJar=%s", s.MinecraftJar())), 0744); err != nil {
			return s.VersionDetails, errors.Wrapf(err, "failed to write fabric properties: %q", fabricSettings)
		}
	}

	// mods
	if s.PurgeMods {
		if err := os.RemoveAll(filepath.Join(s.Path, "mods")); err != nil {
			return s.VersionDetails, errors.Wrap(err, "failed to purge mod directory")
		}
	}

	if err := s.Mods.Install(ctx, s.Path); err != nil {
		return s.VersionDetails, errors.Wrap(err, "failed to install mods")
	}

	return s.VersionDetails, nil
}
