package minecraft

import (
	"context"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/internal/download"
)

type Mods []Mod

type Mod struct {
	Name    string      `hcl:"name,label"`
	URL     string      `hcl:"url"`
	Configs []ModConfig `hcl:"config,block"`
}

type ModConfig struct {
	Path    string `hcl:"path,label"`
	Content string `hcl:"content"`
}

func (m Mods) Install(ctx context.Context, rootPath string) error {
	modPath := filepath.Join(rootPath, "mods")

	if err := os.MkdirAll(modPath, 0744); err != nil {
		return errors.Wrapf(err, "failed to create mods directory %q", modPath)
	}

	for _, mod := range m {
		endpoint, err := url.Parse(mod.URL)
		if err != nil {
			return errors.Wrapf(err, "failed to parse mod %q url %q", mod.Name, mod.URL)
		}

		_, file := path.Split(endpoint.Path)
		modFilePath := filepath.Join(modPath, file)

		if _, err := os.Stat(modFilePath); os.IsNotExist(err) {
			if err := download.File(ctx, mod.URL, modFilePath); err != nil {
				return errors.Wrapf(err, "failed to download mod %q", mod.Name)
			}
		}

		for _, cfg := range mod.Configs {
			configFilePath := filepath.Join(rootPath, cfg.Path)
			configPath, _ := filepath.Split(configFilePath)

			if err := os.MkdirAll(configPath, 0744); err != nil {
				return errors.Wrapf(err, "failed to create mod %q configuration directory %q", mod.Name, configPath)
			}

			if err := os.WriteFile(configFilePath, []byte(cfg.Content), 0744); err != nil {
				return errors.Wrapf(err, "failed to write mod %q config file %q", mod.Name, cfg.Path)
			}
		}
	}

	return nil
}
