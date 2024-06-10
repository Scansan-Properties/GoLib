package ScanSanAPI

import (
	"io"
	l "log/slog"
	"net/http"
)

func webRequest(url string) ([]byte, int, error) {

	// Try 3 Times to get the response

	for i := 0; i < 3; i++ {
		l.With("attempt", i).Debug("Calling ML API")

		resp, err := http.Get(url)
		if err != nil {
			l.With("error", err).Error("Error getting ML Response")
			_ = resp.Body.Close()
			return []byte{}, resp.StatusCode, err
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			l.With("error", err).Error("Error getting ML Response")
			_ = resp.Body.Close()
			return []byte{}, resp.StatusCode, err
		}

		if string(body) != "Internal Server Error" {
			_ = resp.Body.Close()
			return body, resp.StatusCode, nil
		}

		_ = resp.Body.Close()
	}

	return []byte{}, 500, nil
}
