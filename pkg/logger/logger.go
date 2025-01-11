package logger

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

func New(opts *Options) (*slog.Logger, error) {
	const (
		flag     int         = os.O_CREATE | os.O_APPEND | os.O_WRONLY
		dirPerm  fs.FileMode = 0o0755
		filePerm fs.FileMode = 0o0644
	)

    if err := os.MkdirAll(opts.Directory, dirPerm); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

	name := filepath.Join(opts.Directory, "inbox.log")
	f, err := os.OpenFile(name, flag, filePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	w := io.MultiWriter(f, os.Stdout)
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.Level(opts.Level),
	})

	return slog.New(handler), nil
}
