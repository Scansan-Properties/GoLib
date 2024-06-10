package ScanSanAPI

import (
	"io"
	l "log/slog"
	"net/http"
)

func webRequest(url string) ([]byte, int, error) {

	resp, err := http.Get(url)
	if err != nil {
		l.With("error", err).Error("Error getting ML Response")
		return []byte{}, resp.StatusCode, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.With("error", err).Error("Error getting ML Response")
		return []byte{}, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
