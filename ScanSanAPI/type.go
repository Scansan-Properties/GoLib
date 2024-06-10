package ScanSanAPI

import (
	l "log/slog"
	"os"
)

type Request struct {
	Server     string
	DEBUG      bool
	Display    bool
	DebugLevel l.Level
}

func NewRequest() *Request {

	if os.Getenv("SCANSAN_SERVER") == "" {
		l.Error("SCANSAN_SERVER not set")
		return nil
	}

	return &Request{
		Server:  os.Getenv("SCANSAN_SERVER"),
		DEBUG:   false,
		Display: false,
	}
}
