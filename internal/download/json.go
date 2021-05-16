package download

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func JSON(ctx context.Context, file string, v interface{}) error {
	client := newHTTPClient()

	req, err := http.NewRequest("GET", file, nil)
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response")
	}

	if err := json.Unmarshal(data, v); err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}

	return nil
}
