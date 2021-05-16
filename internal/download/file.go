package download

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// File will download the src file to the dest location
func File(ctx context.Context, src, dest string) error {
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

	out, err := os.Create(dest)
	if err != nil {
		return errors.Wrapf(err, "failed to create output file %q", dest)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return errors.Wrapf(err, "failed to download file %q", src)
}
