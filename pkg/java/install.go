package java

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/renevo/mcutils/internal/download"
)

// Install will download and install the versioned JRE into VersionPaths[version]
func (v VersionMap) Install(ctx context.Context, version int, installPath string) error {
	osses, found := v[version]
	if !found {
		return errors.Errorf("unknown java version %d", version)
	}

	architectures, found := osses[runtime.GOOS]
	if !found {
		return errors.Errorf("no version found for %q os", runtime.GOOS)
	}

	endpoint, found := architectures[runtime.GOARCH]
	if !found {
		return errors.Errorf("no version found for cpu architecture %q", runtime.GOARCH)
	}

	outputPath := filepath.Join(installPath, "java", fmt.Sprintf("%d", version))

	if err := os.MkdirAll(outputPath, 0644); err != nil {
		return errors.Wrapf(err, "failed to create java install path %q", outputPath)
	}

	// check for java version already installed
	javaReleaseInfo := filepath.Join(installPath, VersionPaths[version], "release")
	if stat, _ := os.Stat(javaReleaseInfo); stat != nil && stat.Size() > 0 {
		return nil
	}

	if strings.HasSuffix(endpoint, ".zip") {
		return download.ExtractZip(ctx, endpoint, outputPath)
	}

	return download.ExtractTar(ctx, endpoint, outputPath)
}
