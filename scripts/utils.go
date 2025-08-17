package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func getConfigTOMLPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	return filepath.Join(configDir, "hat", "config.toml"), nil
}
