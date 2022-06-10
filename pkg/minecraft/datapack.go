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

type Datapacks []Datapack

type Datapack struct {
	Name string `hcl:"name,label"`
	URL  string `hcl:"url"`
}

func (d Datapacks) Install(ctx context.Context, worldPath string) error {
	datapacksPath := filepath.Join(worldPath, "datapacks")

	if err := os.MkdirAll(datapacksPath, 0744); err != nil {
		return errors.Wrapf(err, "failed to create datapacks directory %q", datapacksPath)
	}

	for _, pack := range d {
		endpoint, err := url.Parse(pack.URL)
		if err != nil {
			return errors.Wrapf(err, "failed to parse mod %q url %q", pack.Name, pack.URL)
		}

		_, file := path.Split(endpoint.Path)
		datapackPath := filepath.Join(datapacksPath, pack.Name+"-"+file)

		if _, err := os.Stat(datapackPath); os.IsNotExist(err) {
			if err := download.File(ctx, pack.URL, datapackPath); err != nil {
				return errors.Wrapf(err, "failed to download mod %q", pack.Name)
			}
		}
	}

	return nil
}
