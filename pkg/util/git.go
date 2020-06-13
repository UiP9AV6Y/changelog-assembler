package util

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// GitUsername returns the name as configured
// in the Git user config. If the config does
// not exist or the key is not set, an empty
// result is given.
func GitUsername() string {
	var err error
	var cfg *ini.File
	var sec *ini.Section
	var key *ini.Key

	path, ok := GitConfig()
	if !ok {
		return ""
	}

	cfg, err = ini.Load(path)
	if err != nil {
		return ""
	}

	sec, err = cfg.GetSection("user")
	if err != nil {
		return ""
	}

	key, err = sec.GetKey("name")
	if err != nil {
		return ""
	}

	return key.String()
}

// Gitconfig returns the path to
// the global Git user config.
func GitConfig() (string, bool) {
	var path string

	path, _ = os.UserConfigDir()
	path = filepath.Join(path, "git", "config")

	if _, err := os.Stat(path); err == nil {
		return path, true
	}

	path, _ = os.UserHomeDir()
	path = filepath.Join(path, ".gitconfig")

	if _, err := os.Stat(path); err == nil {
		return path, true
	}

	return "", false
}
