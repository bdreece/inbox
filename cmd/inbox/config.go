package main

import (
	"os"

	"go.uber.org/config"
)

func loadConfig(paths... string) (*config.YAML, error) {
    opts := make([]config.YAMLOption, len(paths) + 2)
    opts[0] = config.Expand(os.LookupEnv)
	opts[1] = config.Permissive()
    for i, path := range paths {
        opts[i + 2] = config.File(path)
    }

	return config.NewYAML(opts...)
}
