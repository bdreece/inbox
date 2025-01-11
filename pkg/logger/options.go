package logger

import (
	"fmt"

	"go.uber.org/config"
)


type Options struct {
	Level     int    `yaml:"level"`
	Directory string `yaml:"dir"`
}

func Configure(cfg config.Provider) (*Options, error) {
    const key string = "log"
    var opts Options

    if err := cfg.Get(key).Populate(&opts); err != nil {
        return nil, fmt.Errorf("failed to configure logger options: %w", err)
    }

    return &opts, nil
}
