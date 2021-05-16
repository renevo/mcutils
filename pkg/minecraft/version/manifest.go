package version

import (
	"context"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/internal/download"
)

// Manifest of Minecraft versions
type Manifest struct {
	// Latest Release and Snapshot as published by Mojang
	Latest struct {
		// Release is the current version
		Release string `json:"release"`

		// Snapshot is the latest test version of Minecraft
		Snapshot string `json:"snapshot"`
	} `json:"latest"`

	// Versions of all possible Minecraft installs
	Versions []Version `json:"versions"`
}

func GetManifest(ctx context.Context) (Manifest, error) {
	var result Manifest

	if err := download.JSON(ctx, "https://launchermeta.mojang.com/mc/game/version_manifest.json", &result); err != nil {
		return result, errors.Wrap(err, "failed to get launcher metadata")
	}

	return result, nil
}

// GetVersion will lookup and return the specified version. If the version is "latest", it will use what is declared as the latest release, if version is "snapshot" it will lookup what is declared as the latest snapshot.
func (m Manifest) GetVersion(ctx context.Context, version string) (Version, error) {
	if version == "latest" {
		version = m.Latest.Release
	}

	if version == "snapshot" {
		version = m.Latest.Snapshot
	}

	var result Version

	for _, v := range m.Versions {
		if version == v.ID {
			result = v
			break
		}
	}

	if result.ID == "" {
		return result, errors.Errorf("version %q not found", version)
	}

	if err := download.JSON(ctx, result.URL, &result); err != nil {
		return result, errors.Wrapf(err, "failed to download version details for %q", version)
	}

	// set 8 as the default version if not set
	if result.Java.Version == 0 {
		result.Java.Version = 8
	}

	return result, nil
}
