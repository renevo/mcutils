package download

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func ExtractZip(ctx context.Context, src, dest string) error {

	client := newHTTPClient()

	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to get response")
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	// not ideal, but... zip reader isn't as nice as tar/gzip :(
	zipData, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response")
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))

	if err != nil {
		return errors.Wrap(err, "failed to create reader for zip file")
	}

	for _, f := range zipReader.File {
		// Store filename/path for returning and using later on
		target := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
			return errors.Errorf("%s: illegal file path", target)
		}

		if f.FileInfo().IsDir() {
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return errors.Wrapf(err, "failed to create directory %q", target)
				}
			}

			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return errors.Wrapf(err, "failed to create directory %q", filepath.Dir(target))
		}

		// Write File
		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, f.Mode())
		if err != nil {
			return errors.Wrapf(err, "failed to create file %q", target)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return errors.Wrapf(err, "failed to open zipped file %q", f.Name)
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return errors.Wrapf(err, "failed to extract file %q to %q", f.Name, target)
		}

	}

	return nil
}
